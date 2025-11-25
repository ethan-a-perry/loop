// Form
const form = document.getElementById('loop-form');
const formStatus = document.getElementById('form-status');
const pageStatus = document.getElementById('status');
const pageStatusDescription = document.getElementById('status-description');

if (form && formStatus && pageStatus && pageStatusDescription) {
	form.addEventListener('submit', async (e) => {
		e.preventDefault();

		const formData = new FormData(form);

		try {
			const res = await fetch("/api/spotify/loop", {
				method: "POST",
				body: JSON.stringify({
					start: parseInt(formData.get("start")) * 1000,
					end: parseInt(formData.get("end")) * 1000,
				}),
				headers: {
					"Content-Type": "application/json"
				},
			});

			if (res.ok) {
				pageStatus.className = "status looping";
				pageStatusDescription.textContent = "Looping";
			} else {
				formStatus.textContent = 'Apologies, something went wrong when starting the loop.';
			}
		} catch (err) {
			formStatus.textContent = 'Apologies, something went wrong when starting the loop.';
			console.error(err);
		}
	});
}

// Stop Loop Button
const stopLoopBtn = document.getElementById('stop-loop-btn');

if (stopLoopBtn) {
	stopLoopBtn.addEventListener('click', async () => {
		const res = await fetch('/api/spotify/loop/stop');

		if (res.ok) {
			pageStatus.className = "status playing";
			pageStatusDescription.textContent = "Apologies, something went wrong when stopping the loop.";
		}
		else {
			formStatus.textContent = "failure";
		}
	});
}

// Slider
const track = document.querySelector('.track');

const startSlider = document.getElementById('start-slider');
const endSlider = document.getElementById('end-slider');

const startDisplay = document.getElementById('start-display');
const endDisplay = document.getElementById('end-display');

if (startSlider && endSlider && startDisplay && endDisplay) {
	startSlider.addEventListener('input', () => {
		let start = parseInt(startSlider.value, 10);
		const end = parseInt(endSlider.value, 10);

		const minGap = 3;

		if (start + minGap >= end) {
			start = end - minGap;
			startSlider.value = start;
		}

		setDisplay(startDisplay, start);
	})

	endSlider.addEventListener('input', () => {
		const start = parseInt(startSlider.value, 10);
		let end = parseInt(endSlider.value, 10);

		const minGap = 3;

		if (start + minGap >= end) {
			end = start + minGap;
			endSlider.value = end;
		}

		setDisplay(endDisplay, end);
	})
}

function setDisplay(display, value) {
	const minutes = Math.floor(value / 60);
	const seconds = value % 60;
	const timestamp = `${minutes}:${seconds.toString().padStart(2, '0')}`;

	display.textContent = timestamp;
}
