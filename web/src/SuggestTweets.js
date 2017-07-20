import React, { Component } from 'react';
import { Menu } from 'semantic-ui-react'
import { 
  Container,
  Divider,
} from 'semantic-ui-react'
import XHR from './xhr.js';
import TweetCard from './TweetCard.js';
import './App.css';

class SuggestTweets extends Component {
  constructor(props) {
    super(props);

    this.state = {
      tweets: [],
    }

    this.fetch();
  }

  fetch() {
    var params = {
      k: 'golang',
    };
    XHR.getJson(
      XHR.domain + '/api/1.0/suggest',
      params,
    ).then((json) => {
      this.setState({tweets: json});
    }).catch((response) => {
    });
  }

  signin = () => {
    document.location = process.env.REACT_APP_API_DOMAIN + '/api/twitter/signin';
  }

  render() {
    return (
      <Container style={{marginTop: '1em'}}>
        <Menu secondary>
          <Menu.Item name='suggested tweets' active={true} />
          <Menu.Item name='suggested users' />
          <Menu.Menu position='right'>
            <Menu.Item name='sign in' onClick={this.signin} />
          </Menu.Menu>
        </Menu>

        {this.state.tweets.map(
          (tweet) => <div key={tweet.uid}>
            <TweetCard
              name={tweet.name}
              screen_name={tweet.screen_name}
              avatar={tweet.avatar}
              time={tweet.time}
              tweet_id={tweet.tweet_id}
              text={tweet.text}
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
    );
  }
}

export default SuggestTweets;
