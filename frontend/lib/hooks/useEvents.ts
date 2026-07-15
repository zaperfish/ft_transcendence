import { useInfiniteQuery } from "@tanstack/react-query";
import { getMyEvents } from "@/lib/api/events";

const FILTER_MAP = {
	attending: 'member',
	hosting: 'admin',
} as const;

const PAGE_SIZE = 10;

/**
 * Custom hook that fetches paginated events for the current user based on the selected tab.
 * Returns infinite query data for either events the user is attending or hosting.
 */
export function useEvents(tab: 'attending' | 'hosting') {
	const filter = FILTER_MAP[tab];

	return useInfiniteQuery({
		queryKey: ['myEvents', tab],
		queryFn: async ({ pageParam = 1 }) => {
			return getMyEvents({
				filter,
				page: pageParam,
				page_size: PAGE_SIZE,
			});
		},
		getNextPageParam: (lastPage) => {
			const { page, page_size, total } = lastPage;
			const maxPage = Math.ceil(total / page_size);
			return page < maxPage ? page + 1 : undefined;
		},
		initialPageParam: 1,
		networkMode: 'offlineFirst', // in Offline mode reading cache data
	});
}