const isDev = process.env.NODE_ENV === 'development'

/**
 * @type {import('next').NextConfig}
 *
 * Dev mode (`next dev`): keep the Node server so rewrites can proxy /api/* to
 * the Go backend on :3001. Static export is disabled in dev — it also forbids
 * rewrites, so the two must be mutually exclusive.
 *
 * Production build (`next build`): static export for embedding into the Go
 * binary; the backend serves the frontend on the same origin, so no proxy.
 */
const nextConfig = {
  ...(isDev ? {} : { output: 'export' }),
  trailingSlash: true,
}

if (isDev) {
  nextConfig.rewrites = async () => [
    {
      source: '/api/:path*',
      destination: 'http://localhost:3001/api/:path*',
    },
  ]
}

module.exports = nextConfig