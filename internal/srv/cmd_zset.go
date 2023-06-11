package srv

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/weedge/pkg/driver"
	"github.com/weedge/pkg/utils"
)

func init() {
	RegisterCmd("zadd", zadd)
	RegisterCmd("zcard", zcard)
	RegisterCmd("zcount", zcount)
	RegisterCmd("zincrby", zincrby)
	RegisterCmd("zrange", zrange)
	RegisterCmd("zrangebyscore", zrangebyscore)
	RegisterCmd("zrank", zrank)
	RegisterCmd("zrem", zrem)
	RegisterCmd("zremrangebyrank", zremrangebyrank)
	RegisterCmd("zremrangebyscore", zremrangebyscore)
	RegisterCmd("zrevrange", zrevrange)
	RegisterCmd("zrevrank", zrevrank)
	RegisterCmd("zrevrangebyscore", zrevrangebyscore)
	RegisterCmd("zscore", zscore)

	RegisterCmd("zunionstore", zunionstore)
	RegisterCmd("zinterstore", zinterstore)

	RegisterCmd("zrangebylex", zrangebylex)
	RegisterCmd("zremrangebylex", zremrangebylex)
	RegisterCmd("zlexcount", zlexcount)

	// del
	RegisterCmd("zmclear", zmclear)
	// expire
	RegisterCmd("zexpire", zexpire)
	// expireat
	RegisterCmd("zexpireat", zexpireat)
	// ttl
	RegisterCmd("zttl", zttl)
	// persist
	RegisterCmd("zpersist", zpersist)
	// exists
	RegisterCmd("zkeyexists", zkeyexists)
}

func zadd(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 3 {
		err = ErrCmdParams
		return
	}

	args := cmdParams[1:]
	if len(args)%2 == 1 {
		err = ErrCmdParams
		return
	}

	params := make([]driver.ScorePair, len(args)>>1)
	for i := 0; i < len(params); i++ {
		score, err := utils.StrInt64(args[2*i], nil)
		if err != nil {
			return nil, ErrValue
		}

		params[i].Score = score
		params[i].Member = args[2*i+1]
	}

	res, err = c.db.DBZSet().ZAdd(ctx, cmdParams[0], params...)
	return
}

func zcard(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().ZCard(ctx, cmdParams[0])
	return
}

func zcount(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	min, max, err := zparseScoreRange(cmdParams[1], cmdParams[2])
	if err != nil {
		err = ErrValue
		return
	}

	if min > max {
		res = 0
		return
	}

	res, err = c.db.DBZSet().ZCount(ctx, cmdParams[0], min, max)
	return
}

// zparseScoreRange just support int64 score
func zparseScoreRange(minBuf []byte, maxBuf []byte) (min int64, max int64, err error) {
	if strings.ToLower(utils.Bytes2String(minBuf)) == "-inf" {
		min = math.MinInt64
	} else {

		if len(minBuf) == 0 {
			err = ErrCmdParams
			return
		}

		var lopen bool = false
		if minBuf[0] == '(' {
			lopen = true
			minBuf = minBuf[1:]
		}

		min, err = utils.StrInt64(minBuf, nil)
		if err != nil {
			err = ErrValue
			return
		}

		if lopen {
			min++
		}
	}

	if strings.ToLower(utils.Bytes2String(maxBuf)) == "+inf" {
		max = math.MaxInt64
	} else {
		var ropen = false

		if len(maxBuf) == 0 {
			err = ErrCmdParams
			return
		}
		if maxBuf[0] == '(' {
			ropen = true
			maxBuf = maxBuf[1:]
		}

		if maxBuf[0] == '(' {
			ropen = true
			maxBuf = maxBuf[1:]
		}

		max, err = utils.StrInt64(maxBuf, nil)
		if err != nil {
			err = ErrValue
			return
		}

		if ropen {
			max--
		}
	}

	return
}

func zincrby(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	delta, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		return nil, ErrValue
	}

	res, err = c.db.DBZSet().ZIncrBy(ctx, cmdParams[0], delta, cmdParams[2])
	return
}

func zrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return zrangeGeneric(ctx, c, cmdParams, false)
}

func zrangeGeneric(ctx context.Context, c *ConnClient, cmdParams [][]byte, reverse bool) (res interface{}, err error) {
	if len(cmdParams) < 3 {
		return nil, ErrCmdParams
	}

	start, stop, err := zparseRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return nil, ErrValue
	}

	args := cmdParams[3:]
	withScores := false
	if len(args) > 0 {
		if len(args) != 1 {
			return nil, ErrCmdParams
		}
		if strings.ToLower(utils.Bytes2String(args[0])) != "withscores" {
			return nil, ErrSyntax
		}
		withScores = true
	}

	arrScorePair, err := c.db.DBZSet().ZRangeGeneric(ctx, cmdParams[0], start, stop, reverse)
	if err != nil {
		return
	}

	if !withScores {
		members := make([][]byte, 0, len(arrScorePair))
		for _, scorePair := range arrScorePair {
			members = append(members, scorePair.Member)
		}
		res = members
		return
	}

	res = arrScorePair
	return
}

func zparseRange(a1 []byte, a2 []byte) (start int, stop int, err error) {
	if start, err = strconv.Atoi(utils.Bytes2String(a1)); err != nil {
		return
	}

	if stop, err = strconv.Atoi(utils.Bytes2String(a2)); err != nil {
		return
	}

	return
}

func zrangebyscore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return zrangeByScoreGeneric(ctx, c, cmdParams, false)
}

func zrangeByScoreGeneric(ctx context.Context, c *ConnClient, cmdParams [][]byte, reverse bool) (res interface{}, err error) {
	if len(cmdParams) < 3 {
		err = ErrCmdParams
		return
	}

	var minScore, maxScore []byte
	if !reverse {
		minScore, maxScore = cmdParams[1], cmdParams[2]
	} else {
		minScore, maxScore = cmdParams[2], cmdParams[1]
	}

	min, max, err := zparseScoreRange(minScore, maxScore)
	if err != nil {
		return
	}

	args := cmdParams[3:]
	withScores := false
	if len(args) > 0 {
		if strings.ToLower(utils.Bytes2String(args[0])) == "withscores" {
			withScores = true
			args = args[1:]
		}
	}

	offset := 0
	count := -1
	if len(args) > 0 {
		if len(args) != 3 {
			return nil, ErrCmdParams
		}
		if strings.ToLower(utils.Bytes2String(args[0])) != "limit" {
			return nil, ErrSyntax
		}
		if offset, err = strconv.Atoi(utils.Bytes2String(args[1])); err != nil {
			return nil, ErrValue
		}

		if count, err = strconv.Atoi(utils.Bytes2String(args[2])); err != nil {
			return nil, ErrValue
		}
	}

	if offset < 0 {
		return []interface{}{}, nil
	}

	arrScorePair, err := c.db.DBZSet().ZRangeByScoreGeneric(ctx, cmdParams[0], min, max, offset, count, reverse)
	if err != nil {
		return
	}

	if !withScores {
		members := make([][]byte, 0, len(arrScorePair))
		for _, scorePair := range arrScorePair {
			members = append(members, scorePair.Member)
		}
		res = members
		return
	}

	res = arrScorePair
	return
}

func zrank(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	n, err := c.db.DBZSet().ZRank(ctx, cmdParams[0], cmdParams[1])
	if err != nil {
		return
	}
	if n < 0 {
		return
	}

	res = n
	return
}

func zrem(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().ZRem(ctx, cmdParams[0], cmdParams[1:]...)
	return
}

func zremrangebyrank(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	s, e, err := zparseRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return nil, ErrValue
	}

	res, err = c.db.DBZSet().ZRemRangeByRank(ctx, cmdParams[0], s, e)
	return
}

func zremrangebyscore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	s, e, err := zparseRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return nil, ErrValue
	}

	res, err = c.db.DBZSet().ZRemRangeByScore(ctx, cmdParams[0], int64(s), int64(e))
	return
}

func zrevrange(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return zrangeGeneric(ctx, c, cmdParams, true)
}

func zrevrank(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	n, err := c.db.DBZSet().ZRevRank(ctx, cmdParams[0], cmdParams[1])
	if err != nil {
		return
	}
	if n < 0 {
		return
	}

	res = n
	return
}

func zrevrangebyscore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	return zrangeByScoreGeneric(ctx, c, cmdParams, true)
}

func zscore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().ZScore(ctx, cmdParams[0], cmdParams[1])
	if err != nil {
		if err.Error() == "zset score miss" {
			err = nil
		}
		return
	}

	return
}

func zunionstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	destKey, srcKeys, weights, aggregate, err := zparseZsetoptStore(cmdParams)
	if err != nil {
		return
	}

	res, err = c.db.DBZSet().ZUnionStore(ctx, destKey, srcKeys, weights, aggregate)

	return
}

func zparseZsetoptStore(args [][]byte) (destKey []byte, srcKeys [][]byte, weights []int64, aggregate []byte, err error) {
	destKey = args[0]
	nKeys, err := strconv.Atoi(utils.Bytes2String(args[1]))
	if err != nil {
		err = ErrValue
		return
	}
	args = args[2:]
	if len(args) < nKeys {
		err = ErrSyntax
		return
	}

	srcKeys = args[:nKeys]
	args = args[nKeys:]
	weightsFlag := false
	aggregateFlag := false
	for len(args) > 0 {
		op := strings.ToLower(utils.Bytes2String(args[0]))
		switch op {
		case "weights":
			if weightsFlag {
				err = ErrSyntax
				return
			}

			args = args[1:]
			if len(args) < nKeys {
				err = ErrSyntax
				return
			}

			weights = make([]int64, nKeys)
			for i, arg := range args[:nKeys] {
				if weights[i], err = utils.StrInt64(arg, nil); err != nil {
					err = ErrValue
					return
				}
			}
			args = args[nKeys:]
			weightsFlag = true
		case "aggregate":
			if aggregateFlag {
				err = ErrSyntax
				return
			}
			if len(args) < 2 {
				err = ErrSyntax
				return
			}

			op := strings.ToLower(utils.Bytes2String(args[1]))
			switch op {
			case "sum", "min", "max":
				aggregate = args[1]
			default:
				err = ErrSyntax
				return
			}
			args = args[2:]
			aggregateFlag = true
		default:
			err = ErrSyntax
			return
		}
	}
	if !aggregateFlag {
		aggregate = []byte("sum")
	}
	return
}

func zinterstore(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) < 2 {
		err = ErrCmdParams
		return
	}

	destKey, srcKeys, weights, aggregate, err := zparseZsetoptStore(cmdParams)
	if err != nil {
		return
	}

	res, err = c.db.DBZSet().ZInterStore(ctx, destKey, srcKeys, weights, aggregate)
	return
}

func zrangebylex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 && len(cmdParams) != 6 {
		err = ErrCmdParams
		return
	}
	min, max, rangeType, err := zparseMemberRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return
	}

	var offset int = 0
	var count int = -1
	if len(cmdParams) == 6 {
		if strings.ToLower(utils.Bytes2String(cmdParams[3])) != "limit" {
			return nil, ErrSyntax
		}

		if offset, err = strconv.Atoi(utils.Bytes2String(cmdParams[4])); err != nil {
			return nil, ErrValue
		}

		if count, err = strconv.Atoi(utils.Bytes2String(cmdParams[5])); err != nil {
			return nil, ErrValue
		}
	}

	res, err = c.db.DBZSet().ZRangeByLex(ctx, cmdParams[0], min, max, rangeType, offset, count)

	return
}

func zparseMemberRange(minBuf []byte, maxBuf []byte) (min []byte, max []byte, rangeType driver.RangeType, err error) {
	rangeType = driver.RangeClose
	if strings.ToLower(utils.Bytes2String(minBuf)) == "-" {
		min = nil
	} else {
		if len(minBuf) == 0 {
			err = ErrCmdParams
			return
		}

		if minBuf[0] == '(' {
			rangeType |= driver.RangeLOpen
			min = minBuf[1:]
		} else if minBuf[0] == '[' {
			min = minBuf[1:]
		} else {
			err = ErrCmdParams
			return
		}
	}

	if strings.ToLower(utils.Bytes2String(maxBuf)) == "+" {
		max = nil
	} else {
		if len(maxBuf) == 0 {
			err = ErrCmdParams
			return
		}
		if maxBuf[0] == '(' {
			rangeType |= driver.RangeROpen
			max = maxBuf[1:]
		} else if maxBuf[0] == '[' {
			max = maxBuf[1:]
		} else {
			err = ErrCmdParams
			return
		}
	}

	return
}

func zremrangebylex(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	min, max, rangeType, err := zparseMemberRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return nil, err
	}

	res, err = c.db.DBZSet().ZRemRangeByLex(ctx, cmdParams[0], min, max, rangeType)

	return
}

func zlexcount(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 3 {
		err = ErrCmdParams
		return
	}

	min, max, rangeType, err := zparseMemberRange(cmdParams[1], cmdParams[2])
	if err != nil {
		return
	}

	res, err = c.db.DBZSet().ZLexCount(ctx, cmdParams[0], min, max, rangeType)

	return
}

func zmclear(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) == 0 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().Del(ctx, cmdParams...)
	return
}

func zkeyexists(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().Exists(ctx, cmdParams[0])
	return
}

func zexpire(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[2], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBZSet().Expire(ctx, cmdParams[0], d)
	return
}

func zexpireat(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 2 {
		err = ErrCmdParams
		return
	}

	d, err := utils.StrInt64(cmdParams[1], nil)
	if err != nil {
		err = ErrValue
		return
	}

	res, err = c.db.DBZSet().ExpireAt(ctx, cmdParams[0], d)
	return
}

func zttl(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().TTL(ctx, cmdParams[0])
	return
}

func zpersist(ctx context.Context, c *ConnClient, cmdParams [][]byte) (res interface{}, err error) {
	if len(cmdParams) != 1 {
		err = ErrCmdParams
		return
	}

	res, err = c.db.DBZSet().Persist(ctx, cmdParams[0])
	return
}