'use strict';

(async () => {
	const reg = await navigator.serviceWorker.register('sw.js', { scope: "/push" })

	console.log(reg)

	const permission = await Notification.requestPermission()

	console.log(permission)

	// TODO: check if needs to subscribe first

	const pubKeyResponse = await fetch('/pub')

	const pushSubscription = await reg.pushManager.subscribe({
		applicationServerKey: await pubKeyResponse.arrayBuffer(),
		userVisibleOnly: true,
	})


	window.sub = pushSubscription

	console.log(JSON.stringify(pushSubscription, null, "\t"))

	fetch('/subscribe', {
		method: 'POST',
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(pushSubscription),
	})
})()
