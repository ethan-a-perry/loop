const form = document.getElementById('loop-form');
const statusMessage = document.getElementById('form-status');

form.addEventListener('submit', async (e) => {
	e.preventDefault();

	statusMessage.textContent = "Loading...";

	const formData = new FormData(form);

	try {
		const res = await fetch("/api/spotify/loop", {
			method: "POST",
			body: JSON.stringify({
				start: parseInt(formData.get("start")),
				end: parseInt(formData.get("end")),
			}),
			headers: {
				"Content-Type": "application/json"
			},
		});

		if (res.ok) {
			statusMessage.textContent = "Loop is currently active";
		}
		else {
			statusMessage.textContent = "Failed to loop the song";
		}
	} catch (err) {
		statusMessage.textContent = "Failed to loop the song";
		console.error(err);
	}
});
