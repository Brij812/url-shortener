import SignupForm from "./SignupForm";

export const metadata = {
  title: "Signup | HyperLinkOS",
};

export default function SignupPage() {
  return (
    <div className="w-full h-screen flex items-center justify-center">
      <div className="w-full max-w-md p-8 border rounded-lg shadow-sm">
        <h1 className="text-2xl font-semibold mb-6">Create an Account</h1>
        <SignupForm />
      </div>
    </div>
  );
}
