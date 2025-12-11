package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shoeper/gointervals-mcp/internal/config"
	"github.com/shoeper/gointervals-mcp/internal/intervalsclient"
	"io"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetActivities(s *server.MCPServer, config *config.Config) {
    tool := mcp.Tool{
        Name:        "get_activities",
        Description: "Get a list of activities from Intervals.icu",
        InputSchema: mcp.ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                /*"athlete_id": map[string]interface{}{
                    "type":        "string",
                    "description": "The Intervals.icu athlete ID",
                },*/
                "start_date": map[string]interface{}{
                    "type":        "string",
                    "description": "Start date to list activities from in YYYY-MM-DD format",
                    "format": "date",
                },
                "end_date": map[string]interface{}{
                    "type":        "string",
                    "description": "End date to list activities from in YYYY-MM-DD format",
                    "format": "date",
                },
                "limit": map[string]interface{}{
                    "type":        "integer",
                    "default":     10,
                    "maximum": 100,
                    "description": "Maximum number of activities",
                },
            },
            //Required: []string{"athlete_id", "start_date", "end_date"},
        },
    }
    handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        log.Println("get_activities")
        args, ok := request.Params.Arguments.(map[string]interface{})
        if !ok {
            args = make(map[string]interface{})
        }

        // resolve athleteId
        athleteId := config.IntervalsAthleteId
        if athleteId == "" {
            return mcp.NewToolResultError("No athleteID found."), nil
        }

        // resolve apiKey
        apiKey := config.IntervalsApiKey
        if apiKey == "" {
            return mcp.NewToolResultError("No apiKey found."), nil
        }

        // resolve start & end date 
        // default: 30 days ago to today
        now := time.Now()
        endDateStr := now.Format("2006-01-02")
        startDateStr := now.AddDate(0, 0, -30).Format("2006-01-02")

        if val, ok := args["end_date"].(string); ok && val != "" {
            endDateStr = val
        }
        if val, ok := args["start_date"].(string); ok && val != "" {
            startDateStr = val
        }
        startDate, _ := time.Parse("2006-01-02", startDateStr)
        endDate, _ := time.Parse("2006-01-02", endDateStr)

        // use default if startDate is not after endDate
        if startDate.After(endDate) {
            endDateStr = now.Format("2006-01-02")
            startDateStr = now.AddDate(0, 0, -30).Format("2006-01-02")
        }

        // resolve limit
        limit := int32(10)
        if val, ok := args["limit"].(int); ok {
            limit = int32(val)
        }

        params := &intervalsclient.ListActivitiesParams{
            Oldest: startDateStr,
            Newest: &endDateStr,
            Limit: &limit,
        }
        clientOption := intervalsclient.WithRequestEditorFn(intervalsclient.BasicAuth(apiKey))
        client, err := intervalsclient.NewClientWithResponses(config.IntervalsBaseUrl, clientOption)
        if err != nil {
            log.Print("Error creating NewClientWithResponses")
            return mcp.NewToolResultError("Internal error."), nil
        }
        
        resp, err := client.ListActivities(ctx, athleteId, params)
        if err != nil {
            log.Print(err)
            return mcp.NewToolResultError("Error getting activities."), nil
        }
        defer resp.Body.Close()
        if resp.StatusCode != 200 {
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                log.Println("read body error:", err)
            }
            log.Printf("Error getting activities. Status code %d\n%s\n", resp.StatusCode, body)
            return mcp.NewToolResultError(fmt.Sprintf("Error getting activities. Status code %d", resp.StatusCode)), nil
        }
        // Read response body
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                log.Println("read body error:", err)
            }
            log.Printf("Error getting activities. Reading response failed. len: %d\n%s\n", len(body), body)
            return mcp.NewToolResultError("Error getting activities. Reading response failed."), nil
        }
        
        // Parse JSON
        var activities []intervalsclient.Activity  // Define your struct
        if err := json.Unmarshal(body, &activities); err != nil {
            log.Printf("%v", err)
            log.Printf("Error getting activities. Parsing response failed. len: %d\n%s\n", len(body), body)
            return mcp.NewToolResultError("Error getting activities. Parsing response failed."), nil
        }

        var resultText string
        for _, activity := range activities {
            // Parse date and format
            date, err := time.Parse("2006-01-02T15:04:05", *activity.StartDateLocal)
            if err != nil {
                log.Printf("Parsing date failed: %v", err)
            }
            dateStr := date.Format("January 2")
            
            // Format duration (assuming moving_time is in seconds)
            duration := formatDuration(*activity.MovingTime)
            
            // Format distance (assuming distance is in meters)
            distanceKm := *activity.Distance / 1000.0
            
            resultText += fmt.Sprintf("%s: %s - %.1f km in %s, avg HR %d bpm, load: %d\n",
                dateStr,
                *activity.Name,
                distanceKm,
                duration,
                *activity.AverageHeartrate,
                *activity.HrLoad,
            )
        }
        
        return mcp.NewToolResultText(resultText), nil
    }
    
    // Register both the tool definition and its handler
    s.AddTool(tool, handler)
}

func formatDuration(seconds int32) string {
    minutes := seconds / 60
    secs := seconds % 60
    return fmt.Sprintf("%d:%02d", minutes, secs)
}