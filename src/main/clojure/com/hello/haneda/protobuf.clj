(ns com.hello.haneda.protobuf
  (:import
    [com.hello.haneda.api
      Async$Ack
      Async$Ack$Status]))

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
