const connectSpotifyBtn = document.getElementById("connect-spotify");

if (connectSpotifyBtn) {
	connectSpotifyBtn.addEventListener("click", () => {
    window.location.href = "/api/spotify/connect";
	});
}
