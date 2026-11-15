package com.hemingway.kayak.data

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import androidx.room.migration.Migration
import androidx.sqlite.db.SupportSQLiteDatabase
import java.util.Date

/**
 * Base de datos de Room para almacenar modelos de kayak y datos meteorológicos.
 * Room database to store kayak models and weather data.
 */
@Database(
    entities = [KayakModel::class, WeatherData::class, FishingSpot::class],
    version = 3,
    exportSchema = false
)
@TypeConverters(KayakTypeConverters::class, DateConverters::class)
abstract class KayakDatabase : RoomDatabase() {
    abstract fun kayakDao(): KayakDao
    abstract fun weatherDao(): WeatherDao
    abstract fun fishingSpotDao(): FishingSpotDao

    companion object {
        @Volatile
        private var INSTANCE: KayakDatabase? = null

        // Migration from version 1 to 2
        private val MIGRATION_1_2 = object : Migration(1, 2) {
            override fun migrate(database: SupportSQLiteDatabase) {
                database.execSQL(
                    "CREATE TABLE IF NOT EXISTS `weather_data` (" +
                            "`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, " +
                            "`latitude` REAL NOT NULL, " +
                            "`longitude` REAL NOT NULL, " +
                            "`locationName` TEXT NOT NULL, " +
                            "`temperature` REAL NOT NULL, " +
                            "`feelsLike` REAL NOT NULL, " +
                            "`pressure` REAL NOT NULL, " +
                            "`humidity` INTEGER NOT NULL, " +
                            "`windSpeed` REAL NOT NULL, " +
                            "`windDirection` INTEGER NOT NULL, " +
                            "`cloudCover` INTEGER NOT NULL, " +
                            "`precipitation` REAL NOT NULL, " +
                            "`timestamp` INTEGER NOT NULL, " +
                            "`expirationTime` INTEGER NOT NULL)"
                )
            }
        }
        
        // Migration from version 2 to 3
        private val MIGRATION_2_3 = object : Migration(2, 3) {
            override fun migrate(database: SupportSQLiteDatabase) {
                database.execSQL(
                    "CREATE TABLE IF NOT EXISTS `fishing_spots` (" +
                            "`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, " +
                            "`name` TEXT NOT NULL, " +
                            "`description` TEXT NOT NULL, " +
                            "`latitude` REAL NOT NULL, " +
                            "`longitude` REAL NOT NULL, " +
                            "`createdAt` INTEGER NOT NULL, " +
                            "`lastVisited` INTEGER, " +
                            "`rating` INTEGER NOT NULL, " +
                            "`notes` TEXT NOT NULL, " +
                            "`fishCaught` INTEGER NOT NULL, " +
                            "`isShared` INTEGER NOT NULL, " +
                            "`imageUri` TEXT)"
                )
            }
        }

        /**
         * Obtener una instancia de la base de datos.
         * Get a database instance.
         */
        fun getDatabase(context: Context): KayakDatabase {
            return INSTANCE ?: synchronized(this) {
                val instance = Room.databaseBuilder(
                    context.applicationContext,
                    KayakDatabase::class.java,
                    "kayak_database"
                )
                .fallbackToDestructiveMigration()
                .addMigrations(MIGRATION_1_2, MIGRATION_2_3)
                .build()
                INSTANCE = instance
                instance
            }
        }
    }
}

/**
 * Convertidores de tipos para enums en Room.
 * Type converters for enums in Room.
 */
object KayakTypeConverters {
    @androidx.room.TypeConverter
    fun fromKayakType(type: KayakType): String = type.name

    @androidx.room.TypeConverter
    fun toKayakType(typeName: String): KayakType = KayakType.valueOf(typeName)

    @androidx.room.TypeConverter
    fun fromKayakMaterial(material: KayakMaterial): String = material.name

    @androidx.room.TypeConverter
    fun toKayakMaterial(materialName: String): KayakMaterial = KayakMaterial.valueOf(materialName)

    @androidx.room.TypeConverter
    fun fromPropulsionType(propulsion: PropulsionType): String = propulsion.name

    @androidx.room.TypeConverter
    fun toPropulsionType(propulsionName: String): PropulsionType = PropulsionType.valueOf(propulsionName)
}

/**
 * Convertidores para fechas en Room.
 * Date converters for Room.
 */
object DateConverters {
    @androidx.room.TypeConverter
    fun fromTimestamp(value: Long?): Date? {
        return value?.let { Date(it) }
    }

    @androidx.room.TypeConverter
    fun dateToTimestamp(date: Date?): Long? {
        return date?.time
    }
}
