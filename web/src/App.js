import React, { Component } from 'react';
import {
  Container,
  Label,
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
import XHR from './xhr.js';

class App extends Component {
  constructor(props) {
    super(props);

    this.sessionInfos();

    this.state = {
      username: 'Log in',
      loggedin: false,
      avatar: '',
    }
  }

  signin = () => {
    document.location = process.env.REACT_APP_API_DOMAIN + '/api/twitter/signin';
  }

  sessionInfos = () => {
    XHR.getJson(
      XHR.domain + '/api/1.0/session',
    ).then((json) => {
      this.setState({
        username: json.twitter_name,
        avatar: json.twitter_avatar,
        loggedin: true,
      });
    }).catch((response) => {
      // TODO(remy): handle the error
    });
  }

  render() {
    const imageProps = {
      avatar: true,
      spaced: 'right',
      src: this.state.avatar,
    };

    return (
      <Router>
        <Container>
          <Menu pointing secondary>
            <Menu.Item as={NavLink} exact={true} to='/' activeClassName='active'>
              Suggested Tweets
            </Menu.Item>
            <Menu.Item as={NavLink} exact={true} to='/settings' activeClassName='active'>
              Settings
            </Menu.Item>
            {!this.state.loggedin && (
              <Menu.Menu position='right'>
                <Menu.Item name='sign in' onClick={this.signin} />
              </Menu.Menu>
            )}
            {this.state.loggedin && (
              <Menu.Menu position='right'>
                <Label basic content={this.state.username} image={imageProps} />
              </Menu.Menu>
            )}
          </Menu>
          <Route exact path="/" component={SuggestTweets} />
          <Route path="/settings" component={Settings} />
        </Container>
      </Router>
    );
  }
}
export default App;
