import axios from "axios";

const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 10,
});

// Inject the Bearer token on every request if present
api.interceptors.request.use((config) => {
	const token = localStorage.getItem("token");
	if (token) {
		config.headers["Authorization"] = `Bearer ${token}`;
	}
	return config;
});

// Redirect to login on 401 Unauthorized
api.interceptors.response.use(
	(response) => response,
	(error) => {
		if (error.response && error.response.status === 401) {
			localStorage.clear();
			window.location.hash = "#/login";
		}
		return Promise.reject(error);
	}
);

export default api;
