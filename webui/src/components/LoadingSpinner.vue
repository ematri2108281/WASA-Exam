<script setup>
defineProps({
	/** Whether to show the spinner (true) or the slot content (false) */
	loading: {
		type: Boolean,
		required: true,
	},
	/** Optional message shown below the spinner */
	message: {
		type: String,
		default: '',
	},
})
</script>

<template>
	<div v-if="loading" class="spinner-wrapper" role="status" aria-live="polite">
		<div class="spinner-ring" aria-hidden="true" />
		<p v-if="message" class="spinner-message">{{ message }}</p>
		<span class="visually-hidden">{{ message || 'Loading…' }}</span>
	</div>
	<slot v-else />
</template>

<style scoped>
.spinner-wrapper {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 2rem 0;
	gap: 0.75rem;
}

.spinner-ring {
	width: 2.5rem;
	height: 2.5rem;
	border: 3px solid rgba(0, 0, 0, 0.1);
	border-top-color: var(--bs-primary, #0d6efd);
	border-radius: 50%;
	animation: spin 0.7s linear infinite;
}

@keyframes spin {
	to { transform: rotate(360deg); }
}

.spinner-message {
	margin: 0;
	font-size: 0.875rem;
	color: #6c757d;
}
</style>
