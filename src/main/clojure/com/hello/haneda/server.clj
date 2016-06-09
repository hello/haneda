(ns com.hello.haneda.server
  (:require
    [aleph.http :as http]
    [com.hello.haneda
      [message-format :as mf]
      [protobuf :as pb]]
    [compojure.core :as compojure :refer [GET]]
    [compojure.route :as route]
    [manifold [stream :as s] [deferred :as d]])
  (:import
    [com.hello.haneda.api
      Async$Ack
      Async$Ack$Status]))

(defn get-key
  "Fake key store implementation that just calls .getBytes on the senseid string"
  [^String sense-id]
  (.getBytes sense-id))

(defn- request->sense-id
  [request]
  (-> request :headers (get "X-Hello-Sense-Id")))

(defn- response-400
  [body]
  {:status 400
   :headers {"content-type" "application/text"}
   :body body})

(defn dispatch-message-async
  [sense-id key message]
  ;; Just reply with an ack for now after timeout
  ;; TODO return error if message decoding fails
  (let [{:keys [preamble payload-bytes]} (mf/decode-message key message)]
    ;; TODO handle error parsing payload and preamble
    ;; TODO dispatch based on preamble
    ;; Send to relevant service, reply with ack (with ID) and response
    ;; TODO this should actually make a post request.
    ;; The timeout is just to fake an HTTP request right now
    (d/let-flow [response (d/timeout! (d/deferred) 50 nil)
                 ack-response (pb/ack (.getId preamble) (:success pb/ack-status))
                 payload-bytes (.toByteArray ack-response)]
      (mf/encode-message
        key
        (:ack pb/pb-type)
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
