import LoginForm from "./LoginForm";

export const metadata = {
  title: "Login | HyperLinkOS",
};

export default function LoginPage() {
  return (
    <div className="w-full h-screen flex items-center justify-center">
      <div className="w-full max-w-md p-8 border rounded-lg shadow-sm">
        <h1 className="text-2xl font-semibold mb-6">Login</h1>
        <LoginForm />
      </div>
    </div>
  );
}
