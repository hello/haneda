(ns com.hello.haneda.message-format
  (:import
    [com.hello.haneda.api
      Streaming$Preamble
      Streaming$Preamble$auth_type
      Streaming$Preamble$pb_type]
    [java.nio ByteBuffer]
    [javax.crypto Mac]
    [javax.crypto.spec SecretKeySpec]))

(def length-size-bytes 4)
(def hmac-size-bytes 20) ; TODO we'll know this based on the preamble
(def reserved-size-bytes 0)
(def max-chunk-size 2048)

(defn preamble
  ^Streaming$Preamble [^bytes preamble-bytes]
  (Streaming$Preamble/parseFrom preamble-bytes))

;; TODO hmac should be rolling hmac of entire stream up to this point (all chunks)
(defn hmac-sha1
  [^bytes key ^bytes data]
  (let [algorithm "HmacSHA1"
        signing-key (SecretKeySpec. key algorithm)
        mac (Mac/getInstance algorithm)]
    (.init mac signing-key)
    (.doFinal mac data)))

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
          ^bytes payload
          & [id]]
  (let [preamble (cond-> (Streaming$Preamble/newBuilder)
                    :always (.setType pb-type)
                    :always (.setAuth Streaming$Preamble$auth_type/HMAC_SHA1)
                    id (.setId id)
                    :always .build
                    :always .toByteArray)
        hmac (hmac-sha1 key payload)]
    ;; TODO this is single chunk only right now
    (encode-bytes preamble payload hmac)))
