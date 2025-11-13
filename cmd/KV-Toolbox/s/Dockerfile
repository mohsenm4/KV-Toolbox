# مرحله 1: انتخاب پایه با Go و ابزارهای لازم
FROM golang:1.22

# نصب ابزارهای پایه
RUN apt-get update && apt-get install -y openjdk-17-jdk unzip wget git

# تنظیم متغیرهای محیطی برای Android SDK/NDK
ENV ANDROID_HOME /opt/android-sdk
ENV ANDROID_NDK_HOME $ANDROID_HOME/ndk
ENV PATH $PATH:$ANDROID_HOME/cmdline-tools/latest/bin:$ANDROID_HOME/platform-tools:$ANDROID_HOME/emulator:$ANDROID_HOME/ndk

# دانلود و نصب Android SDK Command Line Tools
RUN mkdir -p ${ANDROID_HOME}/cmdline-tools && \
    cd ${ANDROID_HOME}/cmdline-tools && \
    wget https://dl.google.com/android/repository/commandlinetools-linux-11076708_latest.zip -O tools.zip && \
    unzip tools.zip -d latest && \
    rm tools.zip

# نصب SDK و NDK
RUN yes | sdkmanager --sdk_root=${ANDROID_HOME} --licenses
RUN sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0" "ndk;27.0.12077973"

# نصب fyne جدید
RUN go install fyne.io/tools/cmd/fyne@latest

# تنظیم مسیر Go bin
ENV PATH $PATH:/go/bin

# تنظیم پوشه کار
WORKDIR /app

# کپی پروژه داخل کانتینر
COPY . .

# بیلد خروجی اندروید
CMD ["fyne", "package", "-os", "android", "-appID", "com.example.videotomp3", "-icon", "icon.png"]
