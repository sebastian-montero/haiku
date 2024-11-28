export default function NavBar({ onButtonClick }) {
    return (
        <div className="flex flex-col items-center justify-start mb-3">
        <h1 className="text-3xl font-bold text-gray-800 tracking-wide">haiku⿻</h1>

        <div className="flex space-x-1 justify-center mb-5">
        <button
          className="text-sm px-3 py-2 text-black font-bold hover:underline focus:outline-none"
          onClick={() => onButtonClick('sessions')}
        >
          sessions
        </button>
        <button
          className="text-sm px-4 py-2 text-black font-bold hover:underline focus:outline-none"
          onClick={() => onButtonClick('notebooks')}
        >
          notebooks
        </button>
        <button
          className="text-sm px-4 py-2 text-black font-bold hover:underline focus:outline-none"
          onClick={() => onButtonClick('profile')}
        >
          profile
        </button>
            </div>
        </div>
    );
  }