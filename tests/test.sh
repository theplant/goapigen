rm /Users/sunfmin/IdeaProjects/QPTest/src/com/qortex/android/*.java;
go install . && goapigen -pkg=github.com/theplant/qortexapi -lang=java -java_package=com.qortex.android -outdir=/Users/sunfmin/IdeaProjects/QPTest/src
