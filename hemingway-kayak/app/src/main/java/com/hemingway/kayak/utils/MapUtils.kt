package com.hemingway.kayak.utils

import android.content.Context
import android.graphics.Bitmap
import android.graphics.Canvas
import android.os.Bundle
import android.util.Log
import androidx.annotation.DrawableRes
import androidx.core.content.ContextCompat
import org.osmdroid.util.GeoPoint
import org.osmdroid.views.MapView
import org.osmdroid.views.overlay.Marker

object MapUtils {
    private const val TAG = "MapUtils"

    fun saveMapState(mapView: MapView): Bundle {
        val bundle = Bundle()
        val center = mapView.mapCenter
        bundle.putDouble("latitude", center.latitude)
        bundle.putDouble("longitude", center.longitude)
        bundle.putDouble("zoom", mapView.zoomLevelDouble)
        return bundle
    }

    fun restoreMapState(mapView: MapView, bundle: Bundle?) {
        if (bundle == null) {
            Log.d(TAG, "No map state to restore")
            return
        }

        try {
            val latitude = bundle.getDouble("latitude")
            val longitude = bundle.getDouble("longitude")
            val zoom = bundle.getDouble("zoom")

            val geoPoint = GeoPoint(latitude, longitude)
            mapView.controller.animateTo(geoPoint)
            mapView.controller.setZoom(zoom)

            Log.d(TAG, "Map state restored")
        } catch (e: Exception) {
            Log.e(TAG, "Error restoring map state: ${e.message}")
        }
    }

    fun bitmapDescriptorFromVector(context: Context, @DrawableRes vectorResId: Int): Bitmap? {
        return try {
            val vectorDrawable = ContextCompat.getDrawable(context, vectorResId)
            vectorDrawable?.let {
                vectorDrawable.setBounds(0, 0, vectorDrawable.intrinsicWidth, vectorDrawable.intrinsicHeight)
                val bitmap = Bitmap.createBitmap(
                    vectorDrawable.intrinsicWidth,
                    vectorDrawable.intrinsicHeight,
                    Bitmap.Config.ARGB_8888
                )
                val canvas = Canvas(bitmap)
                vectorDrawable.draw(canvas)
                bitmap
            }
        } catch (e: Exception) {
            Log.e(TAG, "Error creating bitmap: ${e.message}")
            null
        }
    }

    fun zoomToShowAllMarkers(mapView: MapView, markers: List<Marker>, padding: Int = 100) {
        if (markers.isEmpty()) {
            Log.d(TAG, "No markers to zoom to")
            return
        }

        try {
            var minLat = Double.MAX_VALUE
            var maxLat = Double.MIN_VALUE
            var minLon = Double.MAX_VALUE
            var maxLon = Double.MIN_VALUE

            markers.forEach { marker ->
                val pos = marker.position
                minLat = minOf(minLat, pos.latitude)
                maxLat = maxOf(maxLat, pos.latitude)
                minLon = minOf(minLon, pos.longitude)
                maxLon = maxOf(maxLon, pos.longitude)
            }

            val centerLat = (minLat + maxLat) / 2
            val centerLon = (minLon + maxLon) / 2
            val center = GeoPoint(centerLat, centerLon)
            
            val latSpan = maxLat - minLat
            val lonSpan = maxLon - minLon
            val span = maxOf(latSpan, lonSpan) * 111000 
            
            mapView.controller.animateTo(center)
            mapView.controller.setZoom(12.0)
            
            Log.d(TAG, "Zoom adjusted to show all markers")
        } catch (e: Exception) {
            Log.e(TAG, "Error zooming to markers: ${e.message}")
        }
    }

    fun addMarkerToMap(
        mapView: MapView,
        context: Context,
        position: GeoPoint,
        title: String,
        snippet: String? = null,
        @DrawableRes iconResId: Int = 0
    ): Marker? {
        return try {
            val marker = Marker(mapView)
            marker.position = position
            marker.title = title
            marker.snippet = snippet
            marker.setAnchor(Marker.ANCHOR_CENTER, Marker.ANCHOR_BOTTOM)
            mapView.overlays.add(marker)
            marker
        } catch (e: Exception) {
            Log.e(TAG, "Error adding marker: ${e.message}")
            null
        }
    }
}
