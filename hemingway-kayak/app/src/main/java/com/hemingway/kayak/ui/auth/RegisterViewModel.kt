package com.hemingway.kayak.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.hemingway.kayak.data.remote.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class RegisterState(
    val loading: Boolean = false,
    val error: String? = null,
    val success: Boolean = false
)

@HiltViewModel
class RegisterViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _state = MutableStateFlow(RegisterState())
    val state: StateFlow<RegisterState> = _state.asStateFlow()

    fun register(name: String, email: String, password: String) {
        val trimmedName = name.trim()
        if (trimmedName.isEmpty() || email.isBlank() || password.isBlank()) {
            _state.update { it.copy(error = "Completá nombre, correo y contraseña") }
            return
        }
        if (password.length < MIN_PASSWORD_LENGTH) {
            _state.update { it.copy(error = "La contraseña tiene que tener al menos 8 caracteres") }
            return
        }
        _state.update { it.copy(loading = true, error = null) }
        viewModelScope.launch {
            val error = authRepository.register(trimmedName, email, password)
            _state.update { it.copy(loading = false, error = error, success = error == null) }
        }
    }

    private companion object {
        const val MIN_PASSWORD_LENGTH = 8
    }
}
