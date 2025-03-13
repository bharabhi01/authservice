'use client';

import ProtectedRoute from '@/components/auth/ProtectedRoute';
import Navbar from '@/components/ui/Navbar';
import { useAuth } from '@/contexts/AuthContext';

export default function DashboardPage() {
    const { user } = useAuth();

    return (
        <ProtectedRoute>
            <div className="min-h-screen bg-gray-50">
                <Navbar />

                <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                    <div className="px-4 py-6 sm:px-0">
                        <div className="bg-white shadow rounded-lg p-6">
                            <h1 className="text-2xl font-bold text-gray-900 mb-4">Dashboard</h1>

                            <div className="border-t border-gray-200 pt-4">
                                <h2 className="text-lg font-medium text-gray-900 mb-2">Welcome to your dashboard</h2>
                                <p className="text-gray-600 mb-4">
                                    You are logged in as <span className="font-medium">{user?.username}</span> with role <span className="font-medium">{user?.role || 'user'}</span>.
                                </p>

                                <div className="bg-gray-50 p-4 rounded-md">
                                    <h3 className="text-md font-medium text-gray-900 mb-2">Your Account Information</h3>
                                    <ul className="space-y-2">
                                        <li className="text-sm text-gray-600">
                                            <span className="font-medium">Username:</span> {user?.username}
                                        </li>
                                        <li className="text-sm text-gray-600">
                                            <span className="font-medium">Email:</span> {user?.email}
                                        </li>
                                        <li className="text-sm text-gray-600">
                                            <span className="font-medium">Name:</span> {user?.first_name} {user?.last_name}
                                        </li>
                                        <li className="text-sm text-gray-600">
                                            <span className="font-medium">Role:</span> {user?.role || 'user'}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        </ProtectedRoute>
    );
}