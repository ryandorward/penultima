import './css/App.css';
import React, {Suspense } from "react";
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'
import { Home } from "./Home"
import Game from "./Game";

function App() {
  return (
    <Router>
			<Switch>				
				<Route path="/game">
					<Suspense fallback={null}>
            <Game />
          </Suspense>
        </Route>
        <Route path="/">
					<Suspense fallback={null}>
						<Home/>
					</Suspense>
				</Route>
      </Switch>      
    </Router>   
  );
}

export default App;