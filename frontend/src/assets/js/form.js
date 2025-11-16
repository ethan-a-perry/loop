const form = document.getElementById('login-form');
const statusMessage = document.getElementById('form-status');

form.addEventListener('submit', async (e) => {
	console.log('trying to sign in')
	e.preventDefault();

	statusMessage.textContent = "Sending...";

	const formData = new FormData(form);

	try {
		const res = await fetch("/api/auth/sign-in/magic-link", {
			method: "POST",
			body: JSON.stringify({
				email: formData.get('email'),
				callbackURL: window.location.origin + "/dashboard",
			}),
      headers: {
      	"Content-Type": "application/json"
      },
    });

		if (res.ok) {
			statusMessage.textContent = "Magic link sent";
		}
		else {
			statusMessage.textContent = "Failed to send magic link";
		}
	} catch (err) {
		statusMessage.textContent = "Failed to send magic link";
		console.error(err);
	}
})
