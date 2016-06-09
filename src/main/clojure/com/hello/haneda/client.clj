(ns com.hello.haneda.client
  (:require
    [aleph.http :as http]
    [com.hello.haneda
      [message-format :as mf]
      [protobuf :as pb]]
    [manifold [stream :as s] [deferred :as d]]))

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
    (partial mf/decode-message key)))

(defn send-message
  [{:keys [conn last-id]} ^bytes key ^bytes request-payload]
  (let [message (mf/encode-message key (:messeji pb/pb-type) request-payload (swap! last-id inc))]
    (s/put! conn message)))
