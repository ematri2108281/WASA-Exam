<script setup>
import { ref } from 'vue'

const props = defineProps({
	/** The error message text to display */
	msg: {
		type: String,
		required: true,
	},
	/** Visual variant: 'danger' | 'warning' | 'info' */
	variant: {
		type: String,
		default: 'danger',
	},
})

const visible = ref(true)
</script>

<template>
	<div
		v-if="visible"
		:class="`alert alert-${props.variant} d-flex justify-content-between align-items-center`"
		role="alert"
	>
		<span>{{ props.msg }}</span>
		<button
			type="button"
			class="btn-close"
			aria-label="Dismiss error"
			@click="visible = false"
		/>
	</div>
</template>

<style scoped>
.alert {
	animation: fadeIn 0.2s ease-in;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(-4px); }
	to   { opacity: 1; transform: translateY(0); }
}
</style>
