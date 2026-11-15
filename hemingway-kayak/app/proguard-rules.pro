# Add project specific ProGuard rules here.
# By default, the flags in this file are appended to flags specified
# in ${sdk.dir}/tools/proguard/proguard-android.txt
# You can edit the include path and order by changing the proguardFiles
# directive in build.gradle.
#
# For more details, see
#   http://developer.android.com/guide/developing/tools/proguard.html

# OkHttp - Start
-dontwarn okhttp3.**
-dontwarn okio.**
-dontwarn javax.annotation.**
-dontwarn org.conscrypt.**
# OkHttp - End

# Retrofit - Start
-keepattributes Signature, InnerClasses, EnclosingMethod
-keepattributes RuntimeVisibleAnnotations, RuntimeVisibleParameterAnnotations
-keepclassmembers,allowshrinking,allowobfuscation interface * {
    @retrofit2.http.* <methods>;
}
-dontwarn org.codehaus.mojo.animal_sniffer.IgnoreJRERequirement
-dontwarn javax.annotation.**
-dontwarn kotlin.Unit
-dontwarn retrofit2.KotlinExtensions
-if interface * { @retrofit2.http.* <methods>; }
-keep,allowobfuscation interface <1>
# Retrofit - End

# Moshi - Start
-keepclasseswithmembers class * {
    @com.squareup.moshi.* <methods>;
}
-keep @com.squareup.moshi.JsonQualifier interface *
-keepclassmembers @com.squareup.moshi.JsonClass class * extends java.lang.Enum {
    <fields>;
    **[] values();
}
# Moshi - End

# Keep Room and its annotations - Start
-keep class androidx.room.** { *; }
-keep @androidx.room.Entity class * { *; }
-keep class androidx.room.EmptyResultSetException { *; }
# Room - End

# Hilt - Start
-keep,allowobfuscation,allowshrinking @dagger.hilt.android.EarlyEntryPoint class *
# Hilt - End

# Prevent proguard from stripping interface information from DataBinding, Hilt, and AndroidX
-keep class androidx.databinding.** { *; }
-keep class androidx.hilt.** { *; }
-keep class androidx.** { *; }
-keep class com.google.android.material.** { *; }

# Keep Application and MainActivity
-keep class com.hemingway.kayak.MainActivity { *; }
-keep class com.hemingway.kayak.KayakApplication { *; }
