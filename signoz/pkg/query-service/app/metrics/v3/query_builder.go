package v3

import (
	"fmt"
	"strings"
	"time"

	"go.signoz.io/signoz/pkg/query-service/constants"
	"go.signoz.io/signoz/pkg/query-service/model"
	v3 "go.signoz.io/signoz/pkg/query-service/model/v3"
	"go.signoz.io/signoz/pkg/query-service/utils"
)

type Options struct {
	PreferRPM bool
}

var aggregateOperatorToPercentile = map[v3.AggregateOperator]float64{
	v3.AggregateOperatorP05:         0.05,
	v3.AggregateOperatorP10:         0.10,
	v3.AggregateOperatorP20:         0.20,
	v3.AggregateOperatorP25:         0.25,
	v3.AggregateOperatorP50:         0.50,
	v3.AggregateOperatorP75:         0.75,
	v3.AggregateOperatorP90:         0.90,
	v3.AggregateOperatorP95:         0.95,
	v3.AggregateOperatorP99:         0.99,
	v3.AggregateOperatorHistQuant50: 0.50,
	v3.AggregateOperatorHistQuant75: 0.75,
	v3.AggregateOperatorHistQuant90: 0.90,
	v3.AggregateOperatorHistQuant95: 0.95,
	v3.AggregateOperatorHistQuant99: 0.99,
}

var aggregateOperatorToSQLFunc = map[v3.AggregateOperator]string{
	v3.AggregateOperatorAvg:     "avg",
	v3.AggregateOperatorMax:     "max",
	v3.AggregateOperatorMin:     "min",
	v3.AggregateOperatorSum:     "sum",
	v3.AggregateOperatorRateSum: "sum",
	v3.AggregateOperatorRateAvg: "avg",
	v3.AggregateOperatorRateMax: "max",
	v3.AggregateOperatorRateMin: "min",
	v3.AggregateOperatorSumRate: "sum",
	v3.AggregateOperatorAvgRate: "avg",
	v3.AggregateOperatorMaxRate: "max",
	v3.AggregateOperatorMinRate: "min",
}

// See https://github.com/SigNoz/signoz/issues/2151#issuecomment-1467249056
var rateWithoutNegative = `If((value - lagInFrame(value, 1, 0) OVER rate_window) < 0, nan, If((ts - lagInFrame(ts, 1, toDate('1970-01-01')) OVER rate_window) >= 86400, nan, (value - lagInFrame(value, 1, 0) OVER rate_window) / (ts - lagInFrame(ts, 1, toDate('1970-01-01')) OVER rate_window))) `

// buildMetricsTimeSeriesFilterQuery builds the sub-query to be used for filtering
// timeseries based on search criteria
func buildMetricsTimeSeriesFilterQuery(fs *v3.FilterSet, groupTags []v3.AttributeKey, mq *v3.BuilderQuery) (string, error) {
	metricName := mq.AggregateAttribute.Key
	aggregateOperator := mq.AggregateOperator
	var conditions []string
	if mq.Temporality == v3.Delta {
		conditions = append(conditions, fmt.Sprintf("metric_name = %s AND temporality = '%s' ", utils.ClickHouseFormattedValue(metricName), v3.Delta))
	} else {
		conditions = append(conditions, fmt.Sprintf("metric_name = %s AND temporality IN ['%s', '%s']", utils.ClickHouseFormattedValue(metricName), v3.Cumulative, v3.Unspecified))
	}

	if fs != nil && len(fs.Items) != 0 {
		for _, item := range fs.Items {
			toFormat := item.Value
			op := v3.FilterOperator(strings.ToLower(strings.TrimSpace(string(item.Operator))))
			// if the received value is an array for like/match op, just take the first value
			// or should we throw an error?
			if op == v3.FilterOperatorLike || op == v3.FilterOperatorRegex || op == v3.FilterOperatorNotLike || op == v3.FilterOperatorNotRegex {
				x, ok := item.Value.([]interface{})
				if ok {
					if len(x) == 0 {
						continue
					}
					toFormat = x[0]
				}
			}

			if op == v3.FilterOperatorContains || op == v3.FilterOperatorNotContains {
				toFormat = fmt.Sprintf("%%%s%%", toFormat)
			}
			fmtVal := utils.ClickHouseFormattedValue(toFormat)
			switch op {
			case v3.FilterOperatorEqual:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') = %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorNotEqual:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') != %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorIn:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') IN %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorNotIn:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') NOT IN %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorLike:
				conditions = append(conditions, fmt.Sprintf("like(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorNotLike:
				conditions = append(conditions, fmt.Sprintf("notLike(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorRegex:
				conditions = append(conditions, fmt.Sprintf("match(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorNotRegex:
				conditions = append(conditions, fmt.Sprintf("not match(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorGreaterThan:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') > %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorGreaterThanOrEq:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') >= %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorLessThan:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') < %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorLessThanOrEq:
				conditions = append(conditions, fmt.Sprintf("JSONExtractString(labels, '%s') <= %s", item.Key.Key, fmtVal))
			case v3.FilterOperatorContains:
				conditions = append(conditions, fmt.Sprintf("like(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorNotContains:
				conditions = append(conditions, fmt.Sprintf("notLike(JSONExtractString(labels, '%s'), %s)", item.Key.Key, fmtVal))
			case v3.FilterOperatorExists:
				conditions = append(conditions, fmt.Sprintf("has(JSONExtractKeys(labels), '%s')", item.Key.Key))
			case v3.FilterOperatorNotExists:
				conditions = append(conditions, fmt.Sprintf("not has(JSONExtractKeys(labels), '%s')", item.Key.Key))
			default:
				return "", fmt.Errorf("unsupported operation")
			}
		}
	}
	queryString := strings.Join(conditions, " AND ")

	var selectLabels string
	if aggregateOperator == v3.AggregateOperatorNoOp || aggregateOperator == v3.AggregateOperatorRate {
		selectLabels = "labels,"
	} else {
		for _, tag := range groupTags {
			selectLabels += fmt.Sprintf(" JSONExtractString(labels, '%s') as %s,", tag.Key, tag.Key)
		}
	}

	filterSubQuery := fmt.Sprintf("SELECT %s fingerprint FROM %s.%s WHERE %s", selectLabels, constants.SIGNOZ_METRIC_DBNAME, constants.SIGNOZ_TIMESERIES_LOCAL_TABLENAME, queryString)

	return filterSubQuery, nil
}

func buildMetricQuery(start, end, step int64, mq *v3.BuilderQuery, tableName string) (string, error) {

	metricQueryGroupBy := mq.GroupBy

	// if the aggregate operator is a histogram quantile, and user has not forgotten
	// the le tag in the group by then add the le tag to the group by
	if mq.AggregateOperator == v3.AggregateOperatorHistQuant50 ||
		mq.AggregateOperator == v3.AggregateOperatorHistQuant75 ||
		mq.AggregateOperator == v3.AggregateOperatorHistQuant90 ||
		mq.AggregateOperator == v3.AggregateOperatorHistQuant95 ||
		mq.AggregateOperator == v3.AggregateOperatorHistQuant99 {
		found := false
		for _, tag := range mq.GroupBy {
			if tag.Key == "le" {
				found = true
				break
			}
		}
		if !found {
			metricQueryGroupBy = append(
				metricQueryGroupBy,
				v3.AttributeKey{
					Key:      "le",
					DataType: v3.AttributeKeyDataTypeString,
					Type:     v3.AttributeKeyTypeTag,
					IsColumn: false,
				},
			)
		}
	}

	filterSubQuery, err := buildMetricsTimeSeriesFilterQuery(mq.Filters, metricQueryGroupBy, mq)
	if err != nil {
		return "", err
	}

	samplesTableTimeFilter := fmt.Sprintf("metric_name = %s AND timestamp_ms >= %d AND timestamp_ms <= %d", utils.ClickHouseFormattedValue(mq.AggregateAttribute.Key), start, end)

	// Select the aggregate value for interval
	queryTmpl :=
		"SELECT %s" +
			" toStartOfInterval(toDateTime(intDiv(timestamp_ms, 1000)), INTERVAL %d SECOND) as ts," +
			" %s as value" +
			" FROM " + constants.SIGNOZ_METRIC_DBNAME + "." + constants.SIGNOZ_SAMPLES_TABLENAME +
			" INNER JOIN" +
			" (%s) as filtered_time_series" +
			" USING fingerprint" +
			" WHERE " + samplesTableTimeFilter +
			" GROUP BY %s" +
			" ORDER BY %s ts"

	// tagsWithoutLe is used to group by all tags except le
	// This is done because we want to group by le only when we are calculating quantile
	// Otherwise, we want to group by all tags except le
	tagsWithoutLe := []string{}
	for _, tag := range mq.GroupBy {
		if tag.Key != "le" {
			tagsWithoutLe = append(tagsWithoutLe, tag.Key)
		}
	}

	groupByWithoutLe := groupBy(tagsWithoutLe...)
	groupTagsWithoutLe := groupSelect(tagsWithoutLe...)
	orderWithoutLe := orderBy(mq.OrderBy, tagsWithoutLe)

	groupBy := groupByAttributeKeyTags(metricQueryGroupBy...)
	groupTags := groupSelectAttributeKeyTags(metricQueryGroupBy...)
	groupSets := groupingSetsByAttributeKeyTags(metricQueryGroupBy...)
	orderBy := orderByAttributeKeyTags(mq.OrderBy, metricQueryGroupBy)

	if len(orderBy) != 0 {
		orderBy += ","
	}
	if len(orderWithoutLe) != 0 {
		orderWithoutLe += ","
	}

	switch mq.AggregateOperator {
	case v3.AggregateOperatorRate:
		// Calculate rate of change of metric for each unique time series
		groupBy = "fingerprint, ts"
		orderBy = "fingerprint, "
		groupTags = "fingerprint,"
		partitionBy := "fingerprint"
		op := "max(value)" // max value should be the closest value for point in time
		subQuery := fmt.Sprintf(
			queryTmpl, "any(labels) as labels, "+groupTags, step, op, filterSubQuery, groupBy, orderBy,
		) // labels will be same so any should be fine
		query := `SELECT %s ts, ` + rateWithoutNegative + ` as value FROM(%s) WINDOW rate_window as (PARTITION BY %s ORDER BY %s ts) `

		query = fmt.Sprintf(query, "labels as fullLabels,", subQuery, partitionBy, orderBy)
		return query, nil
	case v3.AggregateOperatorSumRate, v3.AggregateOperatorAvgRate, v3.AggregateOperatorMaxRate, v3.AggregateOperatorMinRate:
		rateGroupBy := "fingerprint, " + groupBy
		rateGroupTags := "fingerprint, " + groupTags
		rateOrderBy := "fingerprint, " + orderBy
		partitionBy := "fingerprint"
		if len(groupTags) != 0 {
			partitionBy += ", " + groupTags
			partitionBy = strings.Trim(partitionBy, ", ")
		}
		op := "max(value)"
		subQuery := fmt.Sprintf(
			queryTmpl, rateGroupTags, step, op, filterSubQuery, rateGroupBy, rateOrderBy,
		) // labels will be same so any should be fine
		query := `SELECT %s ts, ` + rateWithoutNegative + ` as rate_value FROM(%s) WINDOW rate_window as (PARTITION BY %s ORDER BY %s ts) `
		query = fmt.Sprintf(query, groupTags, subQuery, partitionBy, rateOrderBy)
		query = fmt.Sprintf(`SELECT %s ts, %s(rate_value) as value FROM (%s) WHERE isNaN(rate_value) = 0 GROUP BY %s ORDER BY %s ts`, groupTags, aggregateOperatorToSQLFunc[mq.AggregateOperator], query, groupSets, orderBy)
		return query, nil
	case
		v3.AggregateOperatorRateSum,
		v3.AggregateOperatorRateMax,
		v3.AggregateOperatorRateAvg,
		v3.AggregateOperatorRateMin:
		partitionBy := ""
		if len(groupTags) != 0 {
			partitionBy = "PARTITION BY " + groupTags
			partitionBy = strings.Trim(partitionBy, ", ")
		}
		op := fmt.Sprintf("%s(value)", aggregateOperatorToSQLFunc[mq.AggregateOperator])
		subQuery := fmt.Sprintf(queryTmpl, groupTags, step, op, filterSubQuery, groupSets, orderBy)
		query := `SELECT %s ts, ` + rateWithoutNegative + ` as value FROM(%s) WINDOW rate_window as (%s ORDER BY %s ts) `
		query = fmt.Sprintf(query, groupTags, subQuery, partitionBy, groupTags)
		return query, nil
	case
		v3.AggregateOperatorP05,
		v3.AggregateOperatorP10,
		v3.AggregateOperatorP20,
		v3.AggregateOperatorP25,
		v3.AggregateOperatorP50,
		v3.AggregateOperatorP75,
		v3.AggregateOperatorP90,
		v3.AggregateOperatorP95,
		v3.AggregateOperatorP99:
		op := fmt.Sprintf("quantile(%v)(value)", aggregateOperatorToPercentile[mq.AggregateOperator])
		query := fmt.Sprintf(queryTmpl, groupTags, step, op, filterSubQuery, groupSets, orderBy)
		return query, nil
	case v3.AggregateOperatorHistQuant50, v3.AggregateOperatorHistQuant75, v3.AggregateOperatorHistQuant90, v3.AggregateOperatorHistQuant95, v3.AggregateOperatorHistQuant99:
		rateGroupBy := "fingerprint, " + groupBy
		rateGroupTags := "fingerprint, " + groupTags
		rateOrderBy := "fingerprint, " + orderBy
		partitionBy := "fingerprint"
		if len(groupTags) != 0 {
			partitionBy += ", " + groupTags
			partitionBy = strings.Trim(partitionBy, ", ")
		}
		op := "max(value)"
		subQuery := fmt.Sprintf(
			queryTmpl, rateGroupTags, step, op, filterSubQuery, rateGroupBy, rateOrderBy,
		) // labels will be same so any should be fine
		query := `SELECT %s ts, ` + rateWithoutNegative + ` as rate_value FROM(%s) WINDOW rate_window as (PARTITION BY %s ORDER BY %s ts) `
		query = fmt.Sprintf(query, groupTags, subQuery, partitionBy, rateOrderBy)
		query = fmt.Sprintf(`SELECT %s ts, sum(rate_value) as value FROM (%s) WHERE isNaN(rate_value) = 0 GROUP BY %s ORDER BY %s ts`, groupTags, query, groupSets, orderBy)
		value := aggregateOperatorToPercentile[mq.AggregateOperator]

		query = fmt.Sprintf(`SELECT %s ts, histogramQuantile(arrayMap(x -> toFloat64(x), groupArray(le)), groupArray(value), %.3f) as value FROM (%s) GROUP BY %s ORDER BY %s ts`, groupTagsWithoutLe, value, query, groupByWithoutLe, orderWithoutLe)
		return query, nil
	case v3.AggregateOperatorAvg, v3.AggregateOperatorSum, v3.AggregateOperatorMin, v3.AggregateOperatorMax:
		op := fmt.Sprintf("%s(value)", aggregateOperatorToSQLFunc[mq.AggregateOperator])
		query := fmt.Sprintf(queryTmpl, groupTags, step, op, filterSubQuery, groupSets, orderBy)
		return query, nil
	case v3.AggregateOperatorCount:
		op := "toFloat64(count(*))"
		query := fmt.Sprintf(queryTmpl, groupTags, step, op, filterSubQuery, groupSets, orderBy)
		return query, nil
	case v3.AggregateOperatorCountDistinct:
		op := "toFloat64(count(distinct(value)))"
		query := fmt.Sprintf(queryTmpl, groupTags, step, op, filterSubQuery, groupSets, orderBy)
		return query, nil
	case v3.AggregateOperatorNoOp:
		queryTmpl :=
			"SELECT fingerprint, labels as fullLabels," +
				" toStartOfInterval(toDateTime(intDiv(timestamp_ms, 1000)), INTERVAL %d SECOND) as ts," +
				" any(value) as value" +
				" FROM " + constants.SIGNOZ_METRIC_DBNAME + "." + constants.SIGNOZ_SAMPLES_TABLENAME +
				" INNER JOIN" +
				" (%s) as filtered_time_series" +
				" USING fingerprint" +
				" WHERE " + samplesTableTimeFilter +
				" GROUP BY fingerprint, labels, ts" +
				" ORDER BY fingerprint, labels, ts"
		query := fmt.Sprintf(queryTmpl, step, filterSubQuery)
		return query, nil
	default:
		return "", fmt.Errorf("unsupported aggregate operator")
	}
}

// groupingSets returns a string of comma separated tags for group by clause
// `ts` is always added to the group by clause
func groupingSets(tags ...string) string {
	withTs := append(tags, "ts")
	return fmt.Sprintf(`GROUPING SETS ( (%s), (%s) )`, strings.Join(withTs, ", "), strings.Join(tags, ", "))
}

// groupBy returns a string of comma separated tags for group by clause
// `ts` is always added to the group by clause
func groupBy(tags ...string) string {
	tags = append(tags, "ts")
	return strings.Join(tags, ",")
}

// groupSelect returns a string of comma separated tags for select clause
func groupSelect(tags ...string) string {
	groupTags := strings.Join(tags, ",")
	if len(tags) != 0 {
		groupTags += ", "
	}
	return groupTags
}

func groupingSetsByAttributeKeyTags(tags ...v3.AttributeKey) string {
	groupTags := []string{}
	for _, tag := range tags {
		groupTags = append(groupTags, tag.Key)
	}
	return groupingSets(groupTags...)
}

func groupByAttributeKeyTags(tags ...v3.AttributeKey) string {
	groupTags := []string{}
	for _, tag := range tags {
		groupTags = append(groupTags, tag.Key)
	}
	return groupBy(groupTags...)
}

func groupSelectAttributeKeyTags(tags ...v3.AttributeKey) string {
	groupTags := []string{}
	for _, tag := range tags {
		groupTags = append(groupTags, tag.Key)
	}
	return groupSelect(groupTags...)
}

// orderBy returns a string of comma separated tags for order by clause
// if the order is not specified, it defaults to ASC
func orderBy(items []v3.OrderBy, tags []string) string {
	var orderBy []string
	for _, tag := range tags {
		found := false
		for _, item := range items {
			if item.ColumnName == tag {
				found = true
				orderBy = append(orderBy, fmt.Sprintf("%s %s", item.ColumnName, item.Order))
				break
			}
		}
		if !found {
			orderBy = append(orderBy, fmt.Sprintf("%s ASC", tag))
		}
	}

	return strings.Join(orderBy, ",")
}

func orderByAttributeKeyTags(items []v3.OrderBy, tags []v3.AttributeKey) string {
	var groupTags []string
	for _, tag := range tags {
		groupTags = append(groupTags, tag.Key)
	}
	return orderBy(items, groupTags)
}

func having(items []v3.Having) string {
	var having []string
	for _, item := range items {
		having = append(having, fmt.Sprintf("%s %s %v", "value", item.Operator, utils.ClickHouseFormattedValue(item.Value)))
	}
	return strings.Join(having, " AND ")
}

func reduceQuery(query string, reduceTo v3.ReduceToOperator, aggregateOperator v3.AggregateOperator) (string, error) {
	var selectLabels string
	var groupBy string
	// NOOP and RATE can possibly return multiple time series and reduce should be applied
	// for each uniques series. When the final result contains more than one series we throw
	// an error post DB fetching. Otherwise just return the single data. This is not known until queried so the
	// the query is prepared accordingly.
	if aggregateOperator == v3.AggregateOperatorNoOp || aggregateOperator == v3.AggregateOperatorRate {
		selectLabels = ", any(fullLabels) as fullLabels"
		groupBy = "GROUP BY fingerprint"
	}
	// the timestamp picked is not relevant here since the final value used is show the single
	// chart with just the query value. For the quer
	switch reduceTo {
	case v3.ReduceToOperatorLast:
		query = fmt.Sprintf("SELECT *, timestamp AS ts FROM (SELECT anyLastIf(value, toUnixTimestamp(ts) != 0) as value, anyIf(ts, toUnixTimestamp(ts) != 0) AS timestamp %s FROM (%s) %s)", selectLabels, query, groupBy)
	case v3.ReduceToOperatorSum:
		query = fmt.Sprintf("SELECT *, timestamp AS ts FROM (SELECT sumIf(value, toUnixTimestamp(ts) != 0) as value, anyIf(ts, toUnixTimestamp(ts) != 0) AS timestamp %s FROM (%s) %s)", selectLabels, query, groupBy)
	case v3.ReduceToOperatorAvg:
		query = fmt.Sprintf("SELECT *, timestamp AS ts FROM (SELECT avgIf(value, toUnixTimestamp(ts) != 0) as value, anyIf(ts, toUnixTimestamp(ts) != 0) AS timestamp %s FROM (%s) %s)", selectLabels, query, groupBy)
	case v3.ReduceToOperatorMax:
		query = fmt.Sprintf("SELECT *, timestamp AS ts FROM (SELECT maxIf(value, toUnixTimestamp(ts) != 0) as value, anyIf(ts, toUnixTimestamp(ts) != 0) AS timestamp %s FROM (%s) %s)", selectLabels, query, groupBy)
	case v3.ReduceToOperatorMin:
		query = fmt.Sprintf("SELECT *, timestamp AS ts FROM (SELECT minIf(value, toUnixTimestamp(ts) != 0) as value, anyIf(ts, toUnixTimestamp(ts) != 0) AS timestamp %s FROM (%s) %s)", selectLabels, query, groupBy)
	default:
		return "", fmt.Errorf("unsupported reduce operator")
	}
	return query, nil
}

// PrepareMetricQuery prepares the query to be used for fetching metrics
// from the database
// start and end are in milliseconds
// step is in seconds
func PrepareMetricQuery(start, end int64, queryType v3.QueryType, panelType v3.PanelType, mq *v3.BuilderQuery, options Options) (string, error) {

	// adjust the start and end time to be aligned with the step interval
	start = start - (start % (mq.StepInterval * 1000))
	end = end - (end % (mq.StepInterval * 1000))

	var query string
	var err error
	if mq.Temporality == v3.Delta {
		if panelType == v3.PanelTypeTable {
			query, err = buildDeltaMetricQueryForTable(start, end, mq.StepInterval, mq, constants.SIGNOZ_TIMESERIES_TABLENAME)
		} else {
			query, err = buildDeltaMetricQuery(start, end, mq.StepInterval, mq, constants.SIGNOZ_TIMESERIES_TABLENAME)
		}
	} else {
		if panelType == v3.PanelTypeTable {
			query, err = buildMetricQueryForTable(start, end, mq.StepInterval, mq, constants.SIGNOZ_TIMESERIES_TABLENAME)
		} else {
			query, err = buildMetricQuery(start, end, mq.StepInterval, mq, constants.SIGNOZ_TIMESERIES_TABLENAME)
		}
	}

	if err != nil {
		return "", err
	}

	if options.PreferRPM && (mq.AggregateOperator == v3.AggregateOperatorRate ||
		mq.AggregateOperator == v3.AggregateOperatorSumRate ||
		mq.AggregateOperator == v3.AggregateOperatorAvgRate ||
		mq.AggregateOperator == v3.AggregateOperatorMaxRate ||
		mq.AggregateOperator == v3.AggregateOperatorMinRate ||
		mq.AggregateOperator == v3.AggregateOperatorRateSum ||
		mq.AggregateOperator == v3.AggregateOperatorRateAvg ||
		mq.AggregateOperator == v3.AggregateOperatorRateMax ||
		mq.AggregateOperator == v3.AggregateOperatorRateMin) {
		var selectLabels string
		if mq.AggregateOperator == v3.AggregateOperatorRate {
			selectLabels = "fullLabels,"
		} else {
			selectLabels = groupSelectAttributeKeyTags(mq.GroupBy...)
		}
		query = `SELECT ` + selectLabels + ` ts, ceil(value * 60) as value FROM (` + query + `)`
	}

	if having(mq.Having) != "" {
		query = fmt.Sprintf("SELECT * FROM (%s) HAVING %s", query, having(mq.Having))
	}

	if panelType == v3.PanelTypeValue {
		query, err = reduceQuery(query, mq.ReduceTo, mq.AggregateOperator)
	}
	return query, err
}

func BuildPromQuery(promQuery *v3.PromQuery, step, start, end int64) *model.QueryRangeParams {
	return &model.QueryRangeParams{
		Query: promQuery.Query,
		Start: time.UnixMilli(start),
		End:   time.UnixMilli(end),
		Step:  time.Duration(step * int64(time.Second)),
	}
}
