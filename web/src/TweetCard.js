import React, { Component } from 'react';
import { 
  Container,
  Header,
} from 'semantic-ui-react'

class TweetCard extends Component {
  render() {
    return <Container>
        <Header size='tiny'>
          {this.props.screen_name}
        </Header>
        <p>
          {this.props.text}
        </p>
        <p>
          <a href={this.props.link}>{this.props.link}</a>
        </p>
      </Container>
  }
}

export default TweetCard;
