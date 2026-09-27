import api from './api'
import type { Movie } from '../types/movie'

export async function getMovies(): Promise<Movie[]> {
  const response = await api.get('/movies')
  return response.data
}

export async function getMovieById(movieId: string): Promise<Movie> {
  const response = await api.get(`/movies/${movieId}`)
  return response.data
}

export async function createMovie(movie: Partial<Movie>) {
  const response = await api.post('/movies', movie)
  return response.data
}

export async function updateMovie(
  movieId: string,
  movie: Partial<Movie>
) {
  const response = await api.put(`/movies/${movieId}`, movie)
  return response.data
}

export async function deleteMovie(movieId: string) {
  const response = await api.delete(`/movies/${movieId}`)
  return response.data
}

// import api from './api'
// import type { Movie } from '../types/movie'

// export async function getMovies(): Promise<Movie[]> {
//   const response = await api.get('/movies')

//   console.log('Movies API response:', response.data)

//   // Backend ส่ง Array ตรง ๆ
//   if (Array.isArray(response.data)) {
//     return response.data
//   }

//   // เผื่อ API response ถูกห่อด้วย { data: [...] }
//   if (Array.isArray(response.data?.data)) {
//     return response.data.data
//   }

//   // เผื่อ API response ถูกห่อด้วย { movies: [...] }
//   if (Array.isArray(response.data?.movies)) {
//     return response.data.movies
//   }

//   throw new Error('Invalid movie response format')
// }

// export async function getMovieById(movieId: string): Promise<Movie> {
//   const response = await api.get(`/movies/${movieId}`)

//   return response.data
// }


// import api from './api'
// import type { Movie } from '../types/movie'

// export async function getMovies(): Promise<Movie[]> {
//   const response = await api.get<Movie[]>('/movies')
//   return response.data
// }

// export async function getMovieById(movieId: string): Promise<Movie> {
//   const response = await api.get<Movie>(`/movies/${movieId}`)
//   return response.data
// }

// import api from './api'
// import type { Movie } from '../types/movie'

// export async function getMovies(): Promise<Movie[]> {
//   const response = await api.get('/movies')

//   console.log('GET /movies response:', response.data)

//   if (Array.isArray(response.data)) {
//     return response.data
//   }

//   if (Array.isArray(response.data?.data)) {
//     return response.data.data
//   }

//   if (Array.isArray(response.data?.movies)) {
//     return response.data.movies
//   }

//   throw new Error('Invalid movie response format')
// }

// export async function getMovieById(movieId: string): Promise<Movie> {
//   const response = await api.get(`/movies/${movieId}`)

//   return response.data
// }