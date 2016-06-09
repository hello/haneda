(ns com.hello.haneda.server
  (:require
    [aleph.http :as http]
    [byte-streams :as bs]
    [compojure.core :as compojure :refer [GET]]
    [compojure.route :as route]
    [manifold [stream :as s] [deferred :as d]])
  (:import
    [com.hello.haneda.api
      Streaming$Preamble
      Streaming$Preamble$auth_type
      Streaming$Preamble$pb_type
      Async$Ack
      Async$Ack$Status]
    [java.nio ByteBuffer]
    [javax.crypto Mac]
    [javax.crypto.spec SecretKeySpec]))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;; Server ;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn- request->sense-id
  [request]
  (-> request :headers (get "X-Hello-Sense-Id")))

(defn- response-400
  [body]
  {:status 400
   :headers {"content-type" "application/text"}
   :body body})

(def length-size-bytes 4)
(def hmac-size-bytes 20) ; TODO we'll know this based on the preamble
(def reserved-size-bytes 0)
(def max-chunk-size 2048)

(defn get-key
  [sense-id]
  (.getBytes sense-id))

(defn preamble
  ^Streaming$Preamble [^bytes preamble-bytes]
  (Streaming$Preamble/parseFrom preamble-bytes))

;; TODO hmac should be of entire stream up to this point (all chunks)
(defn hmac-sha1
  [^bytes key ^bytes data]
  (let [algorithm "HmacSHA1"
        signing-key (SecretKeySpec. key algorithm)
        mac (Mac/getInstance algorithm)]
    (.init mac signing-key)
    (.doFinal mac data)))

;; TODO implement rolling HMAC

(defn aslice
  ^bytes [^bytes array from to]
  (java.util.Arrays/copyOfRange array from to))

(defn parse-chunk
  [key ^bytes chunk]
  (let [payload-length (- (alength chunk) hmac-size-bytes)
        payload-chunk (aslice chunk 0 payload-length)
        hmac-chunk (aslice chunk payload-length (alength chunk))
        computed-hmac (hmac-sha1 key payload-chunk)]
    (when-not (java.util.Arrays/equals hmac-chunk computed-hmac)
      (throw (ex-info "bad-hmac" {::error ::bad-hmac})))
    payload-chunk))

(defn decode-length
  [^bytes length-bytes]
  (let [buffer (doto (ByteBuffer/allocate (count length-bytes))
                (.put length-bytes)
                .flip)]
    (.getInt buffer)))

(defn decode-message
  "A message from sense may be chunked.
  First chunk:
  Length | Preamble PB | Length (of total PB) | PB Payload | Reserved | HMAC
  Following chunks:
  PB Payload Continued | Reserved | HMAC

  Returns a map with keys
  :preamble - The Preamble protobuf
  :payload-bytes - The raw bytes of the payload"
  [key ^bytes message]
  (let [message-length (alength message)
        preamble-length (-> message
                          (aslice 0 length-size-bytes)
                          decode-length)
        preamble-bytes (aslice message length-size-bytes (+ length-size-bytes preamble-length))
        preamble (preamble preamble-bytes)
        payload-length-start (+ length-size-bytes preamble-length)
        payload-length-end (+ payload-length-start length-size-bytes)
        payload-length (-> message
                        (aslice payload-length-start payload-length-end)
                        byte-array
                        decode-length)
        chunks (partition-all max-chunk-size message)
        payload-bytes (try
                        (->> (rest chunks)
                          (cons (drop (+ length-size-bytes preamble-length length-size-bytes) (first chunks)))
                          (mapcat (comp (partial parse-chunk key) byte-array))
                          ;; Doall to force the exception to happen here as opposed to later (laziness...)
                          byte-array)
                        (catch clojure.lang.ExceptionInfo ei
                          nil))]
    {:preamble preamble
     :payload-bytes payload-bytes
     :preamble-length preamble-length
     :payload-length payload-length}))

(defn padded-length
  ^bytes [^bytes protobuf]
  (let [length (alength protobuf)
        buffer (ByteBuffer/allocate length-size-bytes)]
    (.putInt buffer length)
    (.array buffer)))

(defn encode-bytes
  ^bytes [^bytes preamble ^bytes payload ^bytes hmac]
  (byte-array
    (concat (padded-length preamble) preamble
            (padded-length payload) payload
            hmac)))

(defn encode-message
  ^bytes [^bytes key
          ^Streaming$Preamble$pb_type pb-type
          ^bytes payload & [id]]
  (let [preamble (cond-> (Streaming$Preamble/newBuilder)
                    :always (.setType pb-type)
                    :always (.setAuth Streaming$Preamble$auth_type/HMAC_SHA1)
                    id (.setId id)
                    :always .build
                    :always .toByteArray)
        hmac (hmac-sha1 key payload)]
    ;; TODO this is single chunk only right now
    (encode-bytes preamble payload hmac)))

(def pb-type
  {:ack Streaming$Preamble$pb_type/ACK
   :batched-periodic-data Streaming$Preamble$pb_type/BATCHED_PERIODIC_DATA
   :sense-log Streaming$Preamble$pb_type/SENSE_LOG
   :sync-response Streaming$Preamble$pb_type/SYNC_RESPONSE
   :matrix-client-message Streaming$Preamble$pb_type/MATRIX_CLIENT_MESSAGE
   :messeji Streaming$Preamble$pb_type/MESSEJI})

(def ack-status
  {:success Async$Ack$Status/SUCCESS
   :client-encoding-error Async$Ack$Status/CLIENT_ENCODING_ERROR
   :client-request-error Async$Ack$Status/CLIENT_REQUEST_ERROR
   :server-error Async$Ack$Status/SERVER_ERROR})

(defn ack
  [message-id ^Async$Ack$Status status]
  (.. (Async$Ack/newBuilder)
    (setMessageId message-id)
    (setStatus status)
    build))

(defn dispatch-message-async
  [sense-id key message]
  ;; Just reply with an ack for now after timeout
  ;; TODO return error if message decoding fails
  (let [{:keys [preamble payload-bytes]} (decode-message key message)]
    ;; TODO handle error parsing payload and preamble
    ;; TODO dispatch based on preamble
    ;; Send to relevant service, reply with ack (with ID) and response
    ;; TODO this should actually make a post request.
    ;; The timeout is just to fake an HTTP request right now
    (d/let-flow [response (d/timeout! (d/deferred) 50 nil)
                 ack-response (ack (.getId preamble) (:success ack-status))
                 payload-bytes (.toByteArray ack-response)]
      (encode-message
        key
        Streaming$Preamble$pb_type/ACK
        payload-bytes))))

(defn dispatch-handler
  [request]
  (prn request)
  (if-let [sense-id (request->sense-id request)]
    (->
      ;; establish the websocket connection
      (d/let-flow [conn (http/websocket-connection request)]
        ;; Take all messages from Sense and ship them to the right place
        (s/consume
          (fn [message]
            (s/connect
              (dispatch-message-async sense-id (get-key sense-id) message)
              conn
              {:downstream? false}))
          conn))
      (d/catch
        (constantly (response-400 "Expected websocket request"))))
    ;; Couldn't get sense id, return 400 response to client
    (response-400 "Could not parse sense ID.")))

(defn handler
  []
  (compojure/routes
    (GET "/dispatch" req (dispatch-handler req))
    (route/not-found "")))

(defn start-server
  []
  (http/start-server (handler) {:port 11000}))

(defn -main
  [& args]
  (let [server (start-server)]
    (println "server started " server)
    server))



;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;; Client ;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defrecord Client [conn last-id]

  java.io.Closeable
  (close [_]
    (.close conn)))

(defn client
  ([sense-id url]
    (map->Client
      {:conn @(http/websocket-client url {:headers {"X-Hello-Sense-Id" sense-id}})
       :last-id (atom 0)}))
  ([sense-id]
    (client sense-id "ws://localhost:11000/dispatch")))

(defn receive-message
  [{:keys [conn]} key]
  (d/chain
    (s/take! conn)
    (partial decode-message key)))

(defn send-message
  [{:keys [conn last-id]} ^bytes key ^bytes request-payload]
  (let [message (encode-message key (:messeji pb-type) request-payload (swap! last-id inc))]
    (s/put! conn message)))
