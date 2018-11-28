let NATS = require('nats')
let nats = NATS.connect({ 'json': true })

nats.subscribe('kube.event.watch', resp => {
  if (resp instanceof NATS.NatsError) {
    console.log(`could not subscribe to kube.watch ${resp}`)
    return
  }
  let slackMsg = {
    channel: 'CDNPXK2KT',
    attachments: [
      {
        fallback: `Namespace: ${resp.Object.involvedObject.namespace} Name: ${resp.Object.involvedObject.name} Kind: ${resp.Object.involvedObject.kind} Reason: ${resp.Object.reason} Message: ${resp.Object.message} Source: ${resp.Object.source.component}`,
        mrkdwn_in: ['text', 'pretext', 'fields'],
        color: '#33A2FF',
        pretext: `:kube: *${resp.cluster}*`,
        text: `*Namespace:* ${resp.Object.involvedObject.namespace}\n*Name:* ${resp.Object.involvedObject.name}\n*Kind:* ${resp.Object.involvedObject.kind}\n*Reason:* ${resp.Object.reason}\n*Message:* ${resp.Object.message}\n*Source:* ${resp.Object.source.component}`,
        footer: `${resp.Object.lastTimestamp}`
      }
    ]
  }
  nats.requestOne('slack.chat.postMessage', slackMsg, {}, 3000, resp => {
    if (resp instanceof NATS.NatsError) {
      console.log(`could not post slack msg ${resp}`)
    }
  })
})
