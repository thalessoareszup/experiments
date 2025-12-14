import { usePlans } from './hooks/usePlans';
import { PlanList } from './components/PlanList';
import './App.css';

function App() {
  const { plans, loading, error, connected } = usePlans();

  return (
    <div className="app">
      <header className="header">
        <h1>Plan Tracker</h1>
        <div className="connection-status">
          <span
            className={`status-dot ${connected ? 'connected' : 'disconnected'}`}
          />
          {connected ? 'Connected' : 'Disconnected'}
        </div>
      </header>

      <main className="main">
        {loading && (
          <div className="loading">Loading plans...</div>
        )}

        {error && (
          <div className="error">
            Error: {error}
          </div>
        )}

        {!loading && !error && (
          <PlanList plans={plans} />
        )}
      </main>

      <footer className="footer">
        <code>plan serve</code> | Real-time plan tracking for AI agents
      </footer>
    </div>
  );
}

export default App;
