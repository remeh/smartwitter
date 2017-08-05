import React, { Component } from 'react';
import {
  Button,
  Checkbox,
  Card,
  Dimmer,
  Header,
  Grid,
  Icon,
  Image,
  Message,
} from 'semantic-ui-react'
import Moment from 'react-moment';
import TweetText from './TweetText.js';
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
      like_count: this.props.like_count,
      retweet_count: this.props.retweet_count,
      retweeting: false,
      retweeted: this.props.retweeted,
      retweetError: '',
      retweetSuccess: '',
      ignored: false,
      ignoring: false,
      ignoreError: '',
      images: this.imagesInTweet(),
    }
  }

  imagesInTweet = () => {
    let rv = [];
    for (let i in this.props.entities) {
      let entity = this.props.entities[i];
      if (entity.type === 'media') {
        rv.push(entity.url);
      }
    }
    return rv;
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
          like_count: this.state.like_count+1,
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
          retweet_count: this.state.retweet_count+1,
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
          ignored: true,
          ignoring: false,
          ignoreError: '',
        });
    }).catch((response) => {
        this.setState({
          ignoring: false,
          ignoreError: 'Either you\'ve already hidden this tweet, either it\'s not available anymore.',
        });
    });
  }

  toggleAutoundo = () => {
    this.setState({
      autoundo: !this.state.autoundo,
    });
  }

  render() {
    return (
      <Dimmer.Dimmable blurring dimmed={this.state.ignored}>
      <Dimmer active={this.state.ignored} inverted>
        <Button onClick={this.props.reload}>This tweet has been hidden. Click here to reload other tweets.</Button>
      </Dimmer>
      <Card fluid>
        <Card.Content>
          <Grid doubling columns="equal">
            <Grid.Column width={12}>
              <Header size='tiny'> <Image src={this.props.avatar} avatar />
                {this.props.name} <span style={{fontSize: '0.8em', color: 'gray'}}>@{this.props.screen_name}</span>
                <a href={this.props.link}><span style={{marginLeft: '1em', fontSize:'0.8em', color: '#999999'}}><Moment fromNow>{this.props.time}</Moment></span></a>
              </Header>
            </Grid.Column>
            <Grid.Column>
              <Button floated='right' basic icon loading={this.state.ignoring} onClick={this.ignore}>
                <Icon name="close" />
              </Button>
            </Grid.Column>
          </Grid>
          <Grid>
            <Grid.Column>
              <TweetText text={this.props.text} entities={this.props.entities} />
            </Grid.Column>
          </Grid>
          <Grid>
            <Grid.Column>
              <Image.Group size='medium'>
              {this.state.images.map(
                (image) => <Image key={Math.random()} src={image} shape="rounded" />
              )}
              </Image.Group>
            </Grid.Column>
          </Grid>
        </Card.Content>
        <Card.Content extra>
          <Grid textAlign='center' verticalAlign='middle' doubling stackable columns="equal">
            <Grid.Column width={4}>
              <Button
                content='Retweet'
                icon='retweet'
                label={''+this.state.retweet_count}
                disabled={this.state.retweeted}
                loading={this.state.retweeting}
                onClick={this.retweet}
              />
            </Grid.Column>
            <Grid.Column width={4}>
              <Button
                content='Like'
                icon='like'
                label={''+this.state.like_count}
                disabled={this.state.liked}
                loading={this.state.liking}
                onClick={this.like}
              />
            </Grid.Column>
            <Grid.Column width={4}>
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
          {this.state.ignoreError && <Message color='red'>
              {this.state.ignoreError}
            </Message>
          }
        </Card.Content>
      </Card>
      </Dimmer.Dimmable>
    )
  }
}

export default TweetCard;
