package com.hemingway.kayak.data.remote

import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TournamentRepository @Inject constructor(
    private val client: BackendClient
) {
    suspend fun tournament(id: String): TournamentDto? =
        runCatching { client.api.tournament(id) }.getOrNull()

    suspend fun tournaments(): List<TournamentDto> =
        runCatching { client.api.tournaments() }.getOrDefault(emptyList())

    suspend fun live(): List<TournamentDto> =
        runCatching { client.api.liveTournaments() }.getOrDefault(emptyList())

    suspend fun ranking(tournamentId: String): List<RankingDto> =
        runCatching { client.api.ranking(tournamentId) }.getOrDefault(emptyList())

    suspend fun participants(tournamentId: String): List<ParticipantDto> =
        runCatching { client.api.participants(tournamentId) }.getOrDefault(emptyList())

    suspend fun catches(tournamentId: String): List<CatchDto> =
        runCatching { client.api.catches(tournamentId) }.getOrDefault(emptyList())

    suspend fun join(tournamentId: String, kayakId: String?): String? = try {
        client.api.joinTournament(tournamentId, JoinRequest(kayakId))
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo unir al torneo"
    }

    suspend fun recordCatch(
        tournamentId: String,
        participantId: String,
        species: String,
        weight: Double,
        latitude: Double?,
        longitude: Double?,
        photoUrl: String?
    ): String? = try {
        client.api.recordCatch(
            tournamentId,
            RecordCatchRequest(participantId, species, weight, latitude, longitude, photoUrl)
        )
        null
    } catch (e: Exception) {
        e.message ?: "No se pudo cargar la captura"
    }
}
