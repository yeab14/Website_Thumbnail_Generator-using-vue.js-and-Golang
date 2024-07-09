<template>
  <div class="container">
    <h1>Website Thumbnail Generator</h1>
    <form @submit.prevent="generateThumbnail">
      <label for="url">Website URL:</label>
      <input type="text" v-model="url" id="url" placeholder="Enter the URL" required />
      <button type="submit" :disabled="loading">
        <span v-if="loading">Generating...</span>
        <span v-else>Generate Thumbnail</span>
      </button>
    </form>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="thumbnailUrl" class="thumbnail">
      <h2>Generated Thumbnail</h2>
      <img :src="thumbnailUrl" alt="Website Thumbnail" />
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      url: '',
      thumbnailUrl: '',
      loading: false,
      error: ''
    };
  },
  methods: {
    async generateThumbnail() {
      this.loading = true;
      this.error = '';
      this.thumbnailUrl = '';
      
      try {
        const response = await fetch('http://localhost:8080/generate', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ url: this.url })
        });
        if (!response.ok) throw new Error('Failed to generate thumbnail');
        const data = await response.json();

        // Construct the full URL to the thumbnail using the base URL of your backend
        this.thumbnailUrl = `http://localhost:8080${data.thumbnailUrl}`;

      } catch (error) {
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.container {
  max-width: 600px;
  margin: 50px auto;
  padding: 20px;
  background: #2e2e2e;
  border-radius: 10px;
  color: #fff;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  font-family: 'Arial', sans-serif;
}

h1 {
  text-align: center;
  margin-bottom: 20px;
  font-size: 2em;
  color: #f9d342;
}

form {
  display: flex;
  flex-direction: column;
}

label {
  font-size: 1.2em;
  margin-bottom: 10px;
}

input {
  padding: 10px;
  font-size: 1em;
  margin-bottom: 20px;
  border: none;
  border-radius: 5px;
}

button {
  padding: 10px;
  font-size: 1.2em;
  background: #f9d342;
  border: none;
  border-radius: 5px;
  color: #2e2e2e;
  cursor: pointer;
  transition: background 0.3s ease;
}

button:disabled {
  background: #ddd;
  cursor: not-allowed;
}

button:hover:not(:disabled) {
  background: #e6bc31;
}

.thumbnail {
  text-align: center;
  margin-top: 20px;
}

img {
  max-width: 100%;
  border: 5px solid #f9d342;
  border-radius: 10px;
}

.error {
  color: #ff4d4d;
  margin-top: 20px;
  text-align: center;
}
</style>
