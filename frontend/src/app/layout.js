import { AuthProvider } from '@/contexts/AuthContext';
import '@/styles/globals.css';

export const metadata = {
  title: 'Auth Service',
  description: 'Centralized Authentication and Authorization Service',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}