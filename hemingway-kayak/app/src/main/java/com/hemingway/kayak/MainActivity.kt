package com.hemingway.kayak

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.app.AppCompatDelegate
import androidx.navigation.fragment.NavHostFragment
import androidx.lifecycle.lifecycleScope
import com.hemingway.kayak.R
import com.hemingway.kayak.data.UserPreferencesManager
import com.hemingway.kayak.data.remote.AuthRepository
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : AppCompatActivity() {

    @Inject
    lateinit var userPreferencesManager: UserPreferencesManager

    @Inject
    lateinit var authRepository: AuthRepository

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        loadUserPreferences()
        applyStartDestination()
    }

    /**
     * Sin sesión arranca en el login; con token va directo al inicio. Se cambia el
     * start destination del grafo antes de que el NavHost lo use.
     */
    private fun applyStartDestination() {
        val navHost = supportFragmentManager.findFragmentById(R.id.nav_host_fragment) as? NavHostFragment
            ?: return
        val navController = navHost.navController
        val graph = navController.navInflater.inflate(R.navigation.nav_graph)
        graph.setStartDestination(if (authRepository.isLoggedIn) R.id.homeFragment else R.id.loginFragment)
        navController.graph = graph
    }

    private fun loadUserPreferences() {
        lifecycleScope.launch {
            val userPreferences = userPreferencesManager.userPreferencesFlow.first()
            applyTheme(userPreferences.darkMode)
        }
    }

    private fun applyTheme(darkMode: Boolean) {
        AppCompatDelegate.setDefaultNightMode(
            if (darkMode) AppCompatDelegate.MODE_NIGHT_YES
            else AppCompatDelegate.MODE_NIGHT_NO
        )
    }
}
