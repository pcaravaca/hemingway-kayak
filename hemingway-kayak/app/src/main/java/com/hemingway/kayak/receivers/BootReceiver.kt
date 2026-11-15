package com.hemingway.kayak.receivers

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.util.Log
import com.hemingway.kayak.service.FishingConditionsService
import dagger.hilt.android.AndroidEntryPoint
import javax.inject.Inject

@AndroidEntryPoint
class BootReceiver : BroadcastReceiver() {

    private val TAG = "BootReceiver"

    override fun onReceive(context: Context?, intent: Intent?) {
        if (intent?.action == Intent.ACTION_BOOT_COMPLETED) {
            Log.d(TAG, "Device booted, starting services")

            context?.let { ctx ->
                val serviceIntent = Intent(ctx, FishingConditionsService::class.java)
                ctx.startForegroundService(serviceIntent)
            }
        }
    }
}
