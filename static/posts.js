let id = null;

function setup() {
	document.querySelector('#edit-post').hidden = true;
	if (document.querySelectorAll('#post-container').length === 0) {
		document.querySelector('#new-post').hidden = false;
		document.querySelector('#new-post').style.visibility = 'visible';
		document.querySelector('#post-adder').hidden = true;
		document.querySelector('#new-post-title').focus();
		document.querySelector('#new-post-title').value = '';
		document.querySelector('#new-post-content').value = '';
	} else {
		document.querySelector('#new-post').hidden = true;
		document.querySelector('#post-adder').hidden = false;
	}
}

document.querySelector('button#post-add').addEventListener('click', () => {
	document.querySelector('div#new-post').hidden = false;
	document.querySelector('div#new-post').style.visibility = 'visibile';
	document.querySelector('#new-post-title').focus();
	document.querySelector('#posts').opacity = 0.6;
});

document.querySelector('button#new-post-cancel').addEventListener('click', () => {
	document.querySelector('div#new-post').hidden = true;
	document.querySelector('#posts').opacity = 1;
});

function registerListeners() {
	for (const element of document.querySelectorAll('button#tool-buttons.delete')) {
		element.addEventListener('click', e => {
			const attribs = e.target.getAttribute('class').split(' ');
			for (const a of attribs) {
				if (a.startsWith('id-')) {
					sendRemove(a.split('-')[1]);
				}
			}
		});
	}

	for (const element of document.querySelectorAll('button#tool-buttons.edit')) {
		element.addEventListener('click', e => {
			const attribs = e.target.getAttribute('class').split(' ');
			for (const a of attribs) {
				if (a.startsWith('id-')) {
					const postId = Number.parseInt(a.split('-')[1]);
					id = postId;
					const content = document.querySelector(`span.id-${id}`).innerText;
					document.querySelector('.edit-post-content').value = content;
					document.querySelector('#edit-post').hidden = false;
					document.querySelector('#edit-post').style.visibility = 'visible';
				}
			}
		});
	}

	for (const element of document.querySelectorAll('img#profile-picture')) {
		element.addEventListener('click', e => {
			const attribs = e.target.getAttribute('class').split(' ');
			for (const a of attribs) {
				if (a.startsWith('user-')) {
					const username = a.split('-')[1];
                    window.location.replace(`/profile?username=${username}`)
				}
			}
		});
	}
}

document.querySelector('#edit-post-cancel').addEventListener('click', e => {
	document.querySelector('#edit-post').hidden = true;
	document.querySelector('#edit-post').style.visibility = 'visible';
});

document.querySelector('#edit-post #edit-post-button').addEventListener('click', e => {
	sendEdit();
});

async function sendEdit() {
	if (id == null) {
		return;
	}

	fetch('/edit', {
		method: 'POST',
		body: JSON.stringify({
			id: Number.parseInt(id),
			content: document.querySelector('textarea.edit-post-content').value,
		}),
		headers: {
			'Content-Type': 'application/json',
		},
	}).then(res => res.text()).then(text => {
		if (JSON.parse(text).redirect == '/') {
			location.reload();
		}
	});
}

async function sendRemove(postId) {
	fetch('/remove', {
		method: 'POST',
		body: JSON.stringify({id: Number.parseInt(postId)}),
		headers: {
			'Content-Type': 'application/json',
		},
	}).then(res => res.text()).then(text => {
		if (JSON.parse(text).redirect == '/') {
			location.reload();
		}
	});
}

setup();
registerListeners();
