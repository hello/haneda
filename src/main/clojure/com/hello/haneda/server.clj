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
      Async$Ack]
    [java.nio ByteBuffer]))

(defn- request->sense-id
  [request]
  (-> request :headers (get "X-Hello-Sense-Id")))

(defn- response-400
  [body]
  {:status 400
   :headers {"content-type" "application/text"}
   :body body})

(def length-size-bytes 4)
(def hmac-size-bytes 40) ; TODO we'll know this based on the preamble
(def reserved-size-bytes 0)
(def max-chunk-size 2048)

(defn get-key
  [sense-id]
  "dummy")

(defn preamble
  ^Streaming$Preamble [^bytes preamble-bytes]
  (Streaming$Preamble/parseFrom preamble-bytes))

(defn hmac
  [key data]
  ;; TODO
  (byte-array (range hmac-size-bytes)))

(defn aslice
  ^bytes [^bytes array from to]
  (java.util.Arrays/copyOfRange array from to))

(defn parse-chunk
  [key ^bytes chunk]
  (let [payload-length (- (alength chunk) hmac-size-bytes)
        payload-chunk (aslice chunk 0 payload-length)
        hmac-chunk (aslice chunk payload-length (alength chunk))
        computed-hmac (hmac key payload-chunk)]
    ;; TODO return error ack
    (when-not (java.util.Arrays/equals hmac-chunk computed-hmac)
      )
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
        ;; TODO first chunk may use different hash function from remaining chunks
        payload-bytes (->> (rest chunks)
                        (cons (drop (+ length-size-bytes preamble-length length-size-bytes) (first chunks)))
                        (mapcat (comp (partial parse-chunk key) byte-array)))]
    {:preamble preamble
     :payload-bytes payload-bytes}))

(defn padded-length
  ^bytes [^bytes protobuf]
  (let [length (alength protobuf)
        buffer (ByteBuffer/allocate length-size-bytes)]
    (.putInt buffer length)
    (.array buffer)))

(defn encode-bytes
  ^bytes [^bytes preamble ^bytes payload ^bytes hmac]
  ;; TODO reserved?
  (byte-array
    (concat (padded-length preamble) preamble
            (padded-length payload) payload
            hmac)))

(defn encode-message
  ^bytes [^Streaming$Preamble$pb_type pb-type ^bytes payload ^bytes hmac]
  (let [preamble (.. (Streaming$Preamble/newBuilder)
                    (setType pb-type)
                    (setAuth Streaming$Preamble$auth_type/HMAC_SHA1)
                    build
                    toByteArray)]
    (encode-bytes preamble payload hmac)))

(defn dispatch-message-async
  [sense-id key message]
  (let [{:keys [preamble payload-bytes]} (decode-message key message)]
    ;; TODO dispatch based on preamble
    ;; Send to relevant service, reply with ack (with ID) and response
    ))

(defn dispatch-handler
  [request]
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
  (http/start-server handler {:port 11000}))

(defn -main
  [& args]
  (let [server (start-server)]
    (println "server started " server)
    server))
