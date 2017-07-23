import React, { Component } from 'react';
import {
  Container,
  Menu,
} from 'semantic-ui-react'
import {
  BrowserRouter as Router,
  NavLink,
  Route,
} from 'react-router-dom'
import './App.css';
import SuggestTweets from './SuggestTweets.js';
import Settings from './Settings.js';

class App extends Component {
  signin = () => {
    document.location = process.env.REACT_APP_API_DOMAIN + '/api/twitter/signin';
  }

  render() {
    return (
      <Router>
        <Container>
          <Menu secondary>
            <Menu.Item as={NavLink} exact={true} to='/' activeClassName='active'>
              Suggested Tweets
            </Menu.Item>
            <Menu.Item as={NavLink} exact={true} to='/settings' activeClassName='active'>
              Settings
            </Menu.Item>
            <Menu.Menu position='right'>
              <Menu.Item name='sign in' onClick={this.signin} />
            </Menu.Menu>
          </Menu>
          <Route exact path="/" component={SuggestTweets} />
          <Route path="/settings" component={Settings} />
        </Container>
      </Router>
    );
  }
}
export default App;
