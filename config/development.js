module.exports = {
  indexer: {
    catchup: {
      enabled: true,
      blocksBatchSize: 1,

      blockNumberRanges: [[221312, 221312]],
    },
    realtime: {
      enabled: false,
    },
  },
};
