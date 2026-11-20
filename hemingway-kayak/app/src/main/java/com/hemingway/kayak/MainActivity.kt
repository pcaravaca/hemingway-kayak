package com.hemingway.kayak

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.appcompat.app.AppCompatDelegate
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.repeatOnLifecycle
import androidx.navigation.NavOptions
import androidx.navigation.fragment.NavHostFragment
import com.google.android.material.bottomnavigation.BottomNavigationView
import com.hemingway.kayak.R
import com.hemingway.kayak.data.UserPreferencesManager
import com.hemingway.kayak.data.remote.AppUpdateManager
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

    @Inject
    lateinit var appUpdateManager: AppUpdateManager

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        loadUserPreferences()
        applyStartDestination()
        appUpdateManager.checkForUpdate()
        observeUpdates()
    }

    /**
     * Avisa si el servidor publica un versionCode más nuevo que el instalado.
     * Si el usuario lo descarta, no vuelve a insistir en la misma sesión.
     */
    private fun observeUpdates() {
        lifecycleScope.launch {
            repeatOnLifecycle(Lifecycle.State.STARTED) {
                appUpdateManager.update.collect { info ->
                    if (info != null) showUpdateDialog(info)
                }
            }
        }
    }

    private fun showUpdateDialog(info: AppUpdateManager.UpdateInfo) {
        val message = "Hay una versión nueva de Hemingway Kayak (${info.versionName})." +
            (info.releaseNotes?.let { "\n\n$it" } ?: "")
        androidx.appcompat.app.AlertDialog.Builder(this)
            .setTitle("Versión nueva disponible")
            .setMessage(message)
            .setPositiveButton("Descargar") { _, _ ->
                try {
                    startActivity(
                        android.content.Intent(
                            android.content.Intent.ACTION_VIEW,
                            android.net.Uri.parse(info.apkUrl)
                        )
                    )
                } catch (e: Exception) {
                    android.widget.Toast.makeText(
                        this, "No se pudo abrir la descarga", android.widget.Toast.LENGTH_SHORT
                    ).show()
                }
                appUpdateManager.dismiss()
            }
            .setNegativeButton("Más tarde") { _, _ -> appUpdateManager.dismiss() }
            .setOnCancelListener { appUpdateManager.dismiss() }
            .show()
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

        // La barra inferior y el grafo comparten los ids de destino. Se navega a mano y
        // no con setupWithNavController: ese helper no dispara el listener cuando el
        // item ya está seleccionado, y desde una pantalla que no es un tab (Spots,
        // detalle de torneo) no se podía volver a Inicio.
        val bottomNav = findViewById<BottomNavigationView>(R.id.bottomNav)
        val startId = graph.startDestinationId
        // Sin restoreState: si no, volver a Inicio re-abre la última pantalla que se
        // hubiera abierto desde ahí (Spots, detalle de torneo) en vez de la pantalla
        // de inicio del tab.
        val options = NavOptions.Builder()
            .setLaunchSingleTop(true)
            .setPopUpTo(startId, inclusive = false)
            .build()
        bottomNav.setOnItemSelectedListener { item ->
            if (navController.currentDestination?.id != item.itemId) {
                navController.navigate(item.itemId, null, options)
            }
            true
        }
        bottomNav.setOnItemReselectedListener { item ->
            // El tab ya es el actual: sólo hace falta limpiar las pantallas apiladas.
            if (navController.currentDestination?.id != item.itemId) {
                navController.navigate(item.itemId, null, options)
            }
            true
        }
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
