import React, { Component } from 'react';
import {
  Button,
  Checkbox,
  Card,
  Header,
  Grid,
  Icon,
  Image,
  Message,
  Statistic,
} from 'semantic-ui-react'
import Moment from 'react-moment';
import XHR from './xhr.js';

class TweetCard extends Component {
  constructor(props) {
    super(props);

    this.state = {
      autoundo: true,
      liking: false,
      liked: this.props.liked,
      likeError: '',
      likeSuccess: '',
      retweeting: false,
      retweeted: this.props.retweeted,
      retweetError: '',
      retweetSuccess: '',
      ignoring: false,
      ignoreError: '',
      ignoreSuccess: '',
    }
  }

  like = (event, button, data) => {
    this.setState({
      liking: true,
    });
    let params = {tid: this.props.tweet_id, au: this.state.autoundo};
    XHR.postJson(
      XHR.domain + '/api/1.0/like',
      params,
    ).then((json) => {
        this.setState({
          liking: false,
          liked: true,
          likeSuccess: 'Successfully liked',
          likeError: '',
        });
    }).catch((response) => {
        this.setState({
          liking: false,
          liked: true,
          likeSuccess: '',
          likeError: 'Either you\'ve already liked this tweet, either it\'s not available anymore.',
        });
    });
  }

  retweet = (event, button, data) => {
    this.setState({
      retweeting: true,
    });
    let params = {tid: this.props.tweet_id, au: this.state.autoundo};
    XHR.postJson(
      XHR.domain + '/api/1.0/retweet',
      params,
    ).then((json) => {
        this.setState({
          retweeting: false,
          retweeted: true,
          retweetSuccess: 'Successfully retweeted',
          retweetError: '',
        });
    }).catch((response) => {
        this.setState({
          retweeting: false,
          retweeted: true,
          retweetSuccess: '',
          retweetError: 'Either you\'ve already retweeted this tweet, either it\'s not available anymore.',
        });
    });
  }

  ignore = (event, button, data) => {
    this.setState({
      ignoring: true,
    });
    let params = {tid: this.props.tweet_id, au: true};
    XHR.postJson(
      XHR.domain + '/api/1.0/ignore',
      params,
    ).then((json) => {
        this.setState({
          ignoring: false,
          ignoreSuccess: 'Successfully ignored',
          ignoreError: '',
        });
    }).catch((response) => {
        this.setState({
          ignoring: false,
          ignoreSuccess: '',
          ignoreError: 'Either you\'ve already ignored this tweet, either it\'s not available anymore.',
        });
    });
  }

  toggleAutoundo = () => {
    this.setState({
      autoundo: !this.state.autoundo,
    });
  }

  render() {
    return <Card fluid padded>
        <Card.Content>
          <Grid doubling columns="equal">
            <Grid.Column width={15}>
              <Header size='tiny'> <Image src={this.props.avatar} avatar />
                {this.props.name} <span style={{fontSize: '0.8em', color: 'gray'}}>@{this.props.screen_name}</span>
                <span style={{marginLeft: '1em', fontSize:'0.8em', color: '#999999'}}><Moment fromNow>{this.props.time}</Moment></span>
              </Header>
            </Grid.Column>
            <Grid.Column>
              <Button basic icon loading={this.state.ignoring} onClick={this.ignore}>
                <Icon name="close" />
              </Button>
            </Grid.Column>
          </Grid>
          <p>
            {this.props.text}
          </p>
          <p>
            <a href={this.props.link}>{this.props.link}</a>
          </p>
        </Card.Content>
        <Card.Content extra>
          <Grid stackable doubling>
            <Grid.Column width={2}>
              <Statistic size='mini' label='Retweets' value={this.props.retweet_count} />
            </Grid.Column>
            <Grid.Column width={2}>
              <Statistic size='mini' label='Likes' value={this.props.like_count} />
            </Grid.Column>
            <Grid.Column width={2}>
              <Button disabled={this.state.liked} loading={this.state.liking} onClick={this.like}>
                Favorite
              </Button>
            </Grid.Column>
            <Grid.Column width={2}>
              <Button disabled={this.state.retweeted} loading={this.state.retweeting} onClick={this.retweet}>
                Retweet
              </Button>
            </Grid.Column>
            <Grid.Column width={8}>
              <Checkbox toggle checked={this.state.autoundo} onClick={this.toggleAutoundo} label="Auto undo in 24h" />
            </Grid.Column>
          </Grid>
          {this.state.likeSuccess && <Message color='green'>
              {this.state.likeSuccess}
            </Message>
          }
          {this.state.likeError && <Message color='red'>
              {this.state.likeError}
            </Message>
          }
          {this.state.retweetSuccess && <Message color='green'>
              {this.state.retweetSuccess}
            </Message>
          }
          {this.state.retweetError && <Message color='red'>
              {this.state.retweetError}
            </Message>
          }
          {this.state.ignoreSuccess && <Message color='green'>
              {this.state.ignoreSuccess}
            </Message>
          }
          {this.state.ignoreError && <Message color='red'>
              {this.state.ignoreError}
            </Message>
          }
        </Card.Content>
      </Card>
  }
}

export default TweetCard;
