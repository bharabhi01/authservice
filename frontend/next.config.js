/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    swcMinify: true,
    // Configure environment variables that should be available to the browser
    env: {
        NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    },
    // Configure image domains if you're using Next.js Image component
    images: {
        domains: ['localhost'],
    },
    // Add rewrites for API proxy in development
    async rewrites() {
        return process.env.NODE_ENV === 'development'
            ? [
                {
                    source: '/api/:path*',
                    destination: 'http://localhost:8080/api/:path*',
                },
            ]
            : [];
    },
};

module.exports = nextConfig;