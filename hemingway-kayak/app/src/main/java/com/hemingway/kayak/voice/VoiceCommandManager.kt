package com.hemingway.kayak.voice

import android.Manifest
import android.app.Activity
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Bundle
import android.speech.RecognitionListener
import android.speech.RecognizerIntent
import android.speech.SpeechRecognizer
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import java.util.*

interface VoiceCommandListener {
    fun onCommandReceived(command: VoiceCommand)
    fun onListeningStarted()
    fun onListeningStopped()
    fun onError(error: String)
}

class VoiceCommandManager(private val activity: Activity) {

    private var speechRecognizer: SpeechRecognizer? = null
    private val isSpanish = Locale.getDefault().language == "es"
    var listener: VoiceCommandListener? = null
    
    private val _isListening = MutableLiveData<Boolean>()
    val isListening: LiveData<Boolean> = _isListening

    private val _lastCommand = MutableLiveData<VoiceCommand>()
    val lastCommand: LiveData<VoiceCommand> = _lastCommand

    fun initialize() {
        if (checkPermissions()) {
            setupSpeechRecognizer()
        }
    }

    fun startListening() {
        if (!checkPermissions()) {
            requestPermissions()
            return
        }

        val intent = Intent(RecognizerIntent.ACTION_RECOGNIZE_SPEECH).apply {
            putExtra(RecognizerIntent.EXTRA_LANGUAGE_MODEL, RecognizerIntent.LANGUAGE_MODEL_FREE_FORM)
            putExtra(RecognizerIntent.EXTRA_LANGUAGE, if (isSpanish) "es-ES" else "en-US")
            putExtra(RecognizerIntent.EXTRA_MAX_RESULTS, 1)
        }

        speechRecognizer?.startListening(intent)
        _isListening.value = true
        listener?.onListeningStarted()
    }

    fun stopListening() {
        speechRecognizer?.stopListening()
        _isListening.value = false
        listener?.onListeningStopped()
    }

    private fun setupSpeechRecognizer() {
        speechRecognizer = SpeechRecognizer.createSpeechRecognizer(activity).apply {
            setRecognitionListener(object : RecognitionListener {
                override fun onReadyForSpeech(params: Bundle?) {}
                override fun onBeginningOfSpeech() {}
                override fun onRmsChanged(rmsdB: Float) {}
                override fun onBufferReceived(buffer: ByteArray?) {}
                override fun onEndOfSpeech() {
                    _isListening.value = false
                }
                override fun onError(error: Int) {
                    _isListening.value = false
                    listener?.onError("Speech recognition error: $error")
                }

                override fun onResults(results: Bundle?) {
                    val matches = results?.getStringArrayList(SpeechRecognizer.RESULTS_RECOGNITION)
                    matches?.firstOrNull()?.let { text ->
                        VoiceCommand.fromRecognizedText(text, isSpanish)?.let { command ->
                            _lastCommand.value = command
                            listener?.onCommandReceived(command)
                        }
                    }
                    startListening()
                }

                override fun onPartialResults(partialResults: Bundle?) {}
                override fun onEvent(eventType: Int, params: Bundle?) {}
            })
        }
    }

    private fun checkPermissions(): Boolean {
        return ContextCompat.checkSelfPermission(
            activity,
            Manifest.permission.RECORD_AUDIO
        ) == PackageManager.PERMISSION_GRANTED
    }

    private fun requestPermissions() {
        ActivityCompat.requestPermissions(
            activity,
            arrayOf(Manifest.permission.RECORD_AUDIO),
            PERMISSION_REQUEST_CODE
        )
    }

    fun destroy() {
        speechRecognizer?.destroy()
        speechRecognizer = null
    }

    companion object {
        private const val PERMISSION_REQUEST_CODE = 123
    }
}
