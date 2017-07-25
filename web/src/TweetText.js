import React, { Component } from 'react';

class TweetText extends Component {
  constructor(props) {
    super(props);

    let text = this.props.text;

    for (let i = 0; i < this.props.entities.length; i++) {
      let entity = props.entities[i];

      switch (entity.type) {
      case 'hashtag':
        text = this.hashtag(text, entity);
        break;
      default:
        break;
      }
    }

    this.state = {
      text: text,
    }
  }

  hashtag = (text, entity) => {
    console.log(text);
    //let before = <span>{text.slice(0, +entity.indices[0])}</span>
    let after = <span>{text.slice(+entity.indices[1], 0)}</span>
    //let ht = <a href={'https://twitter.com'}>{text.substring(+entity.indices[0], +entity.indices[1])}</a>
    return after;
  }

  render() {
    return (
      <p>
        {this.state.text}
      </p>
    );
  }
}

export default TweetText;
