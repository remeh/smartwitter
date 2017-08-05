import React, { Component } from 'react';
import {
  Button,
  Container,
  Divider,
  Dropdown,
  Form,
  Icon,
  Input,
} from 'semantic-ui-react'
import XHR from './xhr.js';

class Settings extends Component {
  constructor(props) {
    super(props);

    this.fetchKeywords();

    this.state = {
      keywords: [],
    }
  }

  fetchKeywords = () => {
    XHR.getJson(
      XHR.domain + '/api/1.0/keywords',
    ).then((json) => {
      let keywords = [];

      for (let i = 0; i < json.length; i++) {
        let k = json[i];
        keywords.push({
          label: k.label,
          keywords: k.keywords,
          position: k.position,
        });
      }

      this.setState({keywords: keywords});
    }).catch((response) => {
      // TODO(remy): error?
    });
  }

  render() {
    return (
      <Container>
        {this.state.keywords.map(
          (keywords) => <KeywordsSettings
                        label={keywords.label}
                        keywords={keywords.keywords}
                        position={keywords.position}
          />
        )}
      </Container>
    )
  }
}

class KeywordsSettings extends Component {
  constructor(props) {
    super(props);

    this.state = {
      label: this.props.label,
      keywords: this.props.keywords,
      position: this.props.position,
    };
  }

  addKeywords = () => {
    let k = this.state.keywords;
    k.push('');
    this.setState({keywords: k});
  }

  removeKeywords = (position) => {
    let k = this.state.keywords;
    k.splice(position, 1);
    this.setState({keywords: k});
  }

  render() {
    const operators = [
      { key: 'with', text: 'With', value: 'with', },
      { key: 'without', text: 'Without', value: 'without', },
    ];

    return (
      <Container>
        <Form>
          <Form.Field>
            <label>Label</label>
            <Input placeholder="Name of these configuration" value={this.state.label} />
          </Form.Field>
          <Form.Field>
            <label>Keywords</label>
            {this.state.keywords.map(
              (keywords, idx) => <Input
                label={<Dropdown defaultValue='with' options={operators} />}
                labelPosition='left'
                placeholder='Counter-Strike'
                value={keywords}
                icon={<Icon name='remove' link onClick={() => { this.removeKeywords(idx); }} />}
              />
            )}
          </Form.Field>
          <Form.Field>
            <Button icon="add" onClick={this.addKeywords} label="Add keyword"/>
          </Form.Field>
        </Form>
      <br />
      <Button primary>Save</Button>
      </Container>
    );
  }
}


export default Settings;
