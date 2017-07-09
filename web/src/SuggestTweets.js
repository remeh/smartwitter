import React, { Component } from 'react';
import { 
  Container,
  Divider,
  Header,
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
      for (let i = 0; i < json.length; i++) {
        json[i].twitter_id = '' + json[i].twitter_id;
        console.log(json[i]);
      }
      this.setState({tweets: json});
    }).catch((response) => {
    });
  }

  render() {
    return (
      <Container style={{marginTop: '1em'}}>
        <Header size="large">
          Suggested Tweets
        </Header>
        {this.state.tweets.map(
          (tweet) => <div key={tweet.uid}>
            <TweetCard
              screen_name={tweet.screen_name}
              text={tweet.text}
              link={tweet.link}
            />
            <Divider />
          </div>
        )}
      </Container>
    );
  }
}

export default SuggestTweets;
