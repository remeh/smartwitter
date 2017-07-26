import React, { Component } from 'react';
import {
  Container,
  Divider,
  Header,
  Input,
} from 'semantic-ui-react'
import XHR from './xhr.js';
import TweetCard from './TweetCard.js';
import './App.css';

class SuggestTweets extends Component {
  constructor(props) {
    super(props);

    this.refreshTimeout = null;

    this.state = {
      tweets: [],
      keywords: 'golang',
      loading: false,
    }

    this.fetch();
  }

  fetch() {
    var params = {
      k: this.state.keywords,
    };
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
        loading: false,
      });
    });
  }

  onChangeKeywords = (event, data) => {
    this.setState({
      keywords: data.value,
      loading: true,
    });
    if (this.refreshTimeout) {
      clearTimeout(this.refreshTimeout);
      this.refreshTimeout = null;
    }
    this.refreshTimeout = setTimeout(() => this.fetch(), 500);
  }

  render() {
    return (
      <Container style={{marginTop: '1em'}}>
        <Container style={{margin: '1em'}}>
          <Header>Keywords</Header>
          <Input loading={this.state.loading} placeholder='golang' onChange={this.onChangeKeywords} />
        </Container>

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
