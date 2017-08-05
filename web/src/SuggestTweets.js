import React, { Component } from 'react';
import {
  Button,
  Container,
  Divider,
  Tab,
} from 'semantic-ui-react'
import XHR from './xhr.js';
import TweetCard from './TweetCard.js';
import './App.css';

class SuggestTweets extends Component {
  constructor(props) {
    super(props);

    this.refreshTimeout = null;

    this.state = {
      tabs: [],
    }

    this.fetchKeywords();
  }

  fetchKeywords() {
    XHR.getJson(
      XHR.domain + '/api/1.0/keywords',
    ).then((json) => {
      let tabs = [];

      for (let i = 0; i < json.length; i++) {
        tabs.push({
          menuItem: json[i].label,
          render: () => <Tweets p={i} />
        });
      }

      tabs.push(this.addTab());

      this.setState({tabs: tabs});
    }).catch((response) => {
      // TODO(remy): error?
    });
  }

  addTab() {
    return {
      menuItem: { key: 'Add', icon: 'add', content: '' },
      render: <div />,
    };
  }

  render() {
    return (
      <Container style={{marginTop: '1em'}}>
        <Tab panes={this.state.tabs} />
      </Container>
    );
  }
}

class Tweets extends Component {
  constructor(props) {
    super(props);

    this.refreshTimeout = null;

    this.state = {
      tweets: [],
      p: this.props.p,
      loading: false,
      reloadDisabled: true,
    }

    setTimeout(() => {
      this.setState({
        reloadDisabled: false,
      });
    }, 5000);

    this.fetch();
  }

  reload = () => {
    this.setState({
      reloadDisabled: true,
    });

    this.fetch();
  }

  fetch = () => {
    var params = {
      p: this.state.p,
    };

    setTimeout(() => {
      this.setState({
        reloadDisabled: false,
      });
    }, 5000);

    XHR.getJson(
      XHR.domain + '/api/1.0/suggest',
      params,
    ).then((json) => {
      this.setState({
        tweets: json,
        loading: false,
      });
    }).catch((response) => {
      this.setState({
        reloadDisabled: false,
        loading: false,
      });
    });
  }

  render() {
    return (
        <Tab.Pane>
        <Button disabled={this.state.reloadDisabled} onClick={this.reload} icon='refresh' primary content='Reload' />
        <Button icon='edit' content='Configure keywords' />
        <Divider />
        <Container>
          {this.state.tweets.map(
          (tweet) => <div key={tweet.uid}>
            <TweetCard
              reload={this.fetch}
              name={tweet.name}
              screen_name={tweet.screen_name}
              avatar={tweet.avatar}
              time={tweet.time}
              tweet_id={tweet.tweet_id}
              text={tweet.text}
              entities={tweet.entities}
              link={tweet.link}
              retweeted={tweet.retweeted}
              liked={tweet.liked}
              like_count={tweet.like_count}
              retweet_count={tweet.retweet_count}
            />
            <Divider />
          </div>
        )}
        </Container>
      </Tab.Pane>
    );
  }
};

export default SuggestTweets;
