const message = document.querySelector('#status-message');

function checkStatus() {
	if (message.innerText) {
		return;
	}

	message.style.animation = 'none';
}

document.querySelector('button.delete-profile').addEventListener('click', async () => {
	const res = await fetch('/profile', {
		method: 'DELETE',
	});
	const d = JSON.parse(await res.text());
	message.style.animation = 'none';
	message.innerText = d.reason + '. goodbye';
	message.offsetWidth;
	message.style.animation = '';
	setTimeout(() => {
		globalThis.location.href = '/';
	}, 1200);
});

document.querySelector('button.update-picture').addEventListener('click', () => {
	document.querySelector('input.picture').click();
});

document.querySelector('input.picture').addEventListener('change', () => {
	document.querySelector('input.submit-picture').click();
});

checkStatus();
