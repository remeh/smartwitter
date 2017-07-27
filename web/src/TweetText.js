import React, { Component } from 'react';

class TweetText extends Component {
  constructor(props) {
    super(props);

    let tags = this.format(this.props.text);

    this.state = {
      content: tags,
    }
  }

  format = (text) => {

    // no entities? directly return the text.
    if (this.props.entities.length === 0) {
      return [<span>text</span>];
    }

    let tags = [];
    let idx = 0;
    for (let i = 0; i < this.props.entities.length; i++) {
      let entity = this.props.entities[i];
      if (entity.indices[0] !== idx) {
        let substr = text.substring(idx, entity.indices[0]);
        tags.push(<span key={Math.random()}>{substr}</span>);
        idx = entity.indices[0];
      }

      switch (entity.type) {
      case 'hashtag':
        tags.push(this.hashtag(this.props.text, entity));
        idx = entity.indices[1];
        break;
      case 'url':
        tags.push(this.url(this.props.text, entity));
        idx = entity.indices[1];
        break;
      case 'user_mention':
        tags.push(this.userMention(this.props.text, entity));
        idx = entity.indices[1];
      default:
        break;
      }
    }

    // end of text without entity
    // ----------------------
    let lastEnt = this.props.entities[this.props.entities.length-1];
    if (lastEnt.indices[1] != text.length) {
      let substr = text.substring(lastEnt.indices[1]);
      tags.push(<span key={Math.random()}>{substr}</span>);
    }

    return tags;
  }

  hashtag = (text, entity) => {
    let ht = text.substring(+entity.indices[0], +entity.indices[1]);
    let url = 'https://twitter.com/'+ht;
    return <span key={Math.random()}><a href={url}>{ht}</a> </span>
  }

  url = (text, entity) => {
    let t = text.substring(+entity.indices[0], +entity.indices[1]);
    return <span key={Math.random()}><a href={entity.url}>{entity.display_url}</a></span>
  }

  userMention = (text, entity) => {
    let t = text.substring(+entity.indices[0], +entity.indices[1]);
    let url = 'https://twitter.com/'+entity.screen_name;
    return <span key={Math.random()}><a href={url}>@{entity.screen_name}</a></span>
  }

  render() {
    return (
      <div>
        {this.state.content}
      </div>
    );
  }
}

export default TweetText;
