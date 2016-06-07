(defproject com.hello/haneda "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.7.0"]
                 [aleph "0.4.1-beta2"]
                 [com.google.protobuf/protobuf-java "2.6.1"]
                 [compojure "1.4.0"]]
  :source-paths ["src" "src/main/clojure"]
  ;; Java source is stored separately.
  :java-source-paths ["src/main/java"])
